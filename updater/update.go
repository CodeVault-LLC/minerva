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
)

func StartAutoUpdate(interval time.Duration) {
	for _, list := range constants.VC.Lists {
		parsedData, err := fetchAndParseList(&list)
		if err != nil {
			log.Printf("Failed to update %s: %v", list.Description, err)
			continue
		}
		log.Printf("Updated %s with %d entries", list.Description, len(parsedData))

		var data []string
		for item := range parsedData {
			data = append(data, item.Value)
		}

		err = storeParsedData(list.ListID, data)
		if err != nil {
			log.Printf("Failed to store data for %s: %v", list.Description, err)
		}
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		for _, list := range constants.VC.Lists {
			parsedData, err := fetchAndParseList(&list)
			if err != nil {
				log.Printf("Failed to update %s: %v", list.Description, err)
				continue
			}
			log.Printf("Updated %s with %d entries", list.Description, len(parsedData))

			var data []string
			for item := range parsedData {
				data = append(data, item.Value)
			}

			err = storeParsedData(list.ListID, data)
			if err != nil {
				log.Printf("Failed to store data for %s: %v", list.Description, err)
			}
		}
	}
}

func fetchAndParseList(list *config.List) (<-chan parsers.Item, error) {
	resp, err := http.Get(list.URL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println("Parsing data for", list.Description, list.Parser)

	parsedData, err := parsers.ParseBytes(data, list.Parser)
	if err != nil {
		return nil, err
	}

	return parsedData, nil
}

// Updated to store data in DragonflyDB (Redis-compatible)
func storeParsedData(listID string, parsedData []string) error {
	// Use Redis' pipeline to insert multiple records efficiently
	pipe := constants.Rdb.Pipeline()

	for _, value := range parsedData {
		// Store each value in Redis with the ListID as the key
		pipe.SAdd(constants.Ctx, listID, value)
	}

	_, err := pipe.Exec(constants.Ctx)
	return err
}

func CompareValues(comparedValue string) []config.List {
	var lists []config.List

	for _, list := range constants.VC.Lists {
		for _, cat := range list.Categories {
			if cat == comparedValue {
				lists = append(lists, config.List{
					Description: list.Description,
					ListID:      list.ListID,
					Categories:  list.Categories,
					URL:         list.URL,
				})
			}
		}
	}

	return lists
}
