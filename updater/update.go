package updater

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/codevault-llc/humblebrag-api/config"
	"github.com/codevault-llc/humblebrag-api/constants"
	"github.com/codevault-llc/humblebrag-api/parsers"
	"github.com/codevault-llc/humblebrag-api/types"
	"github.com/go-redis/redis/v8"
)

const (
	batchSize          = 5000 // Reduced batch size
	maxRetries         = 5    // Increased retries
	retryDelay         = 10 * time.Second
	storeTimeout       = 5 * time.Minute // Increased timeout
	updateWorkers      = 5
	redisScriptTimeout = 10 * time.Minute // Increased timeout
)

var (
	storeLuaScript = redis.NewScript(`
    local key = KEYS[1]
    local values = ARGV
    local batchSize = 1000  -- Process 1000 items at a time within Redis
    local totalAdded = 0

    redis.call('DEL', key)

    for i = 1, #values, batchSize do
        local batch = {}
        for j = i, math.min(i + batchSize - 1, #values) do
            table.insert(batch, values[j])
        end
        local added = redis.call('SADD', key, unpack(batch))
        totalAdded = totalAdded + added
    end

    redis.call('EXPIRE', key, 1800)  -- 30 minutes
    return totalAdded
`)
)

func StartAutoUpdate(interval time.Duration) {
	updateChan := make(chan *types.List, len(config.ConfigLists))
	var wg sync.WaitGroup

	// Start worker goroutines
	for i := 0; i < updateWorkers; i++ {
		wg.Add(1)
		go updateWorker(updateChan, &wg)
	}

	// Initial update
	for _, list := range config.ConfigLists {
		updateChan <- list
	}

	// Periodic updates
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		for _, list := range config.ConfigLists {
			updateChan <- list
		}
	}

	close(updateChan)
	wg.Wait()
}

func updateWorker(updateChan <-chan *types.List, wg *sync.WaitGroup) {
	defer wg.Done()

	for list := range updateChan {
		if err := updateList(list); err != nil {
			log.Printf("Failed to update %s: %v", list.Description, err)
		}
	}
}

func updateList(list *types.List) error {
	parsedData, err := fetchAndParseList(list)
	if err != nil {
		return fmt.Errorf("failed to fetch and parse list: %w", err)
	}

	log.Printf("Parsed %s with %d entries", list.Description, len(parsedData))

	storedCount, err := storeParsedDataWithRetry(list.ListID, parsedData)
	if err != nil {
		return fmt.Errorf("failed to store data: %w", err)
	}

	log.Printf("Successfully stored %d/%d entries for %s", storedCount, len(parsedData), list.Description)

	if storedCount != len(parsedData) {
		log.Printf("Warning: Mismatch in stored data count for %s. Expected: %d, Actual: %d", list.Description, len(parsedData), storedCount)
	}

	return nil
}

func fetchAndParseList(list *types.List) ([]parsers.Item, error) {
	resp, err := http.Get(list.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch list: %w", err)
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	parsedDataChan, err := parsers.ParseBytes(bodyBytes, list.Parser)
	if err != nil {
		return nil, fmt.Errorf("failed to parse data: %w", err)
	}

	var parsedData []parsers.Item
	for item := range parsedDataChan {
		parsedData = append(parsedData, item)
	}

	return parsedData, nil
}

func storeParsedDataWithRetry(listID string, parsedData []parsers.Item) (int, error) {
	var totalStored int
	for attempt := 0; attempt < maxRetries; attempt++ {
		storedCount, err := storeParsedData(listID, parsedData)
		if err != nil {
			log.Printf("Attempt %d failed to store data for %s: %v", attempt+1, listID, err)
			time.Sleep(retryDelay)
			continue
		}
		totalStored += storedCount
		if totalStored == len(parsedData) {
			return totalStored, nil
		}
		log.Printf("Attempt %d: Stored %d/%d items for %s", attempt+1, totalStored, len(parsedData), listID)
	}
	return totalStored, fmt.Errorf("failed to store all data after %d attempts", maxRetries)
}

func storeParsedData(listID string, parsedData []parsers.Item) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), storeTimeout)
	defer cancel()

	totalStored := 0
	itemsByType := make(map[string][]string)
	for _, item := range parsedData {
		key := fmt.Sprintf("%s:%s", listID, item.Type)
		itemsByType[key] = append(itemsByType[key], item.Value)
	}

	for key, values := range itemsByType {
		storedCount, err := storeDataBatch(ctx, key, values)
		if err != nil {
			return totalStored, fmt.Errorf("failed to store batch for key %s: %w", key, err)
		}
		totalStored += storedCount
	}

	return totalStored, nil
}

func storeDataBatch(ctx context.Context, key string, batch []string) (int, error) {
	result, err := storeLuaScript.Run(ctx, constants.Rdb, []string{key}, batch).Result()
	if err != nil {
		return 0, err
	}
	storedCount, ok := result.(int64)
	if !ok {
		return 0, fmt.Errorf("unexpected result type from Lua script")
	}
	return int(storedCount), nil
}

func CompareValues(comparedValue string, valueType parsers.ListType) []types.List {
	var matchingLists []types.List
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pipe := constants.Rdb.Pipeline()
	cmds := make(map[string]*redis.BoolCmd)

	for _, list := range config.ConfigLists {
		for _, listType := range list.Types {
			if listType == valueType {
				key := fmt.Sprintf("%s:%s", list.ListID, valueType)
				fmt.Println(key)
				cmds[list.ListID] = pipe.SIsMember(ctx, key, comparedValue)
			}
		}
	}

	_, err := pipe.Exec(ctx)
	if err != nil {
		log.Printf("Failed to execute pipeline: %v", err)
		return nil
	}

	for _, list := range config.ConfigLists {
		if cmd, ok := cmds[list.ListID]; ok {
			exists, err := cmd.Result()
			if err != nil {
				log.Printf("Failed to get result for %s: %v", list.ListID, err)
				continue
			}
			if exists {
				matchingLists = append(matchingLists, *list)
			}
		}
	}

	return matchingLists
}
