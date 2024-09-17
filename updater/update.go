package updater

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/codevault-llc/humblebrag-api/config"
	"github.com/codevault-llc/humblebrag-api/constants"
	"github.com/codevault-llc/humblebrag-api/parsers"
	"github.com/codevault-llc/humblebrag-api/types"
)

func StartAutoUpdate(interval time.Duration) {
	for _, list := range config.ConfigLists {
		parsedData, err := fetchAndParseList(list)
		if err != nil {
			log.Printf("Failed to update %s: %v", list.Description, err)
			continue
		}
		log.Printf("Updated %s with %d entries", list.Description, len(parsedData))

		err = storeParsedData(list.ListID, parsedData)
		if err != nil {
			log.Printf("Failed to store data for %s: %v", list.Description, err)
		}
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		for _, list := range config.ConfigLists {
			parsedData, err := fetchAndParseList(list)
			if err != nil {
				log.Printf("Failed to update %s: %v", list.Description, err)
				continue
			}

			err = storeParsedData(list.ListID, parsedData)
			if err != nil {
				log.Printf("Failed to store data for %s: %v", list.Description, err)
			}
		}
	}
}

func fetchAndParseList(list *types.List) ([]parsers.Item, error) {
	resp, err := http.Get(list.URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	parsedDataChan, err := parsers.ParseBytes(data, list.Parser)
	if err != nil {
		fmt.Println("Failed to parse data for", list.Description, err)
		return nil, err
	}

	// Collect items from the channel
	var parsedData []parsers.Item
	for item := range parsedDataChan {
		parsedData = append(parsedData, item)
	}

	return parsedData, nil
}

func storeParsedData(listID string, parsedData []parsers.Item) error {
	// Use Redis' pipeline to insert multiple records efficiently
	pipe := constants.Rdb.Pipeline()

	for _, item := range parsedData {
		// Create a composite key with ListID and item type
		key := fmt.Sprintf("%s:%s", listID, item.Type)

		// Store each value in Redis with the ListID and type as key
		pipe.SAdd(constants.Ctx, key, item.Value)

		// Set the expiration time (30 minutes) for the data
		pipe.Expire(constants.Ctx, key, 30*time.Minute)
	}

	_, err := pipe.Exec(constants.Ctx)
	return err
}

// CompareValues function to search for a value in DragonflyDB and return matching lists
func CompareValues(comparedValue string, valueType parsers.ListType) []types.List {
	var matchingLists []types.List

	// Iterate over all lists
	for _, list := range config.ConfigLists {
		for _, listType := range list.Types {
			// Only search in lists that match the given type
			if listType == valueType {
				// Create the composite key with ListID and type
				key := fmt.Sprintf("%s:%s", list.ListID, valueType)

				// Use SISMEMBER to check if the compared value exists in the Redis set for the list
				exists, err := constants.Rdb.SIsMember(constants.Ctx, key, comparedValue).Result()
				if err != nil {
					log.Printf("Failed to search in Redis for %s: %v", list.ListID, err)
					continue
				}

				// If the value exists in the set, add the list to the result
				if exists {
					matchingLists = append(matchingLists, *list)
				}
			}
		}
	}

	return matchingLists
}
