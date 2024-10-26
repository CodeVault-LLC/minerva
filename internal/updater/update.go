package updater

import (
	"fmt"
	"sync"
	"time"

	"github.com/codevault-llc/humblebrag-api/config"
	"github.com/codevault-llc/humblebrag-api/pkg/types"
)

const (
	updateWorkers = 5
)

func StartAutoUpdate(interval time.Duration) {
	updateChan := make(chan *types.Filter, len(config.ConfigLists))
	var wg sync.WaitGroup

	// Start worker goroutines
	for i := 0; i < updateWorkers; i++ {
		wg.Add(1)
		go updateWorker(updateChan, &wg) // Uncomment this line
	}

	for _, list := range config.ConfigLists {
		updateChan <- list
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		fmt.Println("Updating list")

		for _, list := range config.ConfigLists {
			updateChan <- list
		}
	}

	close(updateChan)
	wg.Wait()
}
