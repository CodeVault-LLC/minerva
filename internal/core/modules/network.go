package modules

import (
	"fmt"
	"sync"

	"github.com/codevault-llc/humblebrag-api/internal/core/modules/network"
	"github.com/codevault-llc/humblebrag-api/internal/models/entities"
)

// NetworkModule orchestrates network-related scans through sub-modules
type NetworkModule struct {
	modules []MiniModule
}

func NewNetworkModule() *NetworkModule {
	return &NetworkModule{
		modules: []MiniModule{
			&network.IPLookupModule{},
			&network.IPRangeLookupModule{},
		},
	}
}

// Execute runs the Network-specific scan logic
func (m *NetworkModule) Execute(job entities.JobModel) error {
	results := make(map[string]interface{})
	errChan := make(chan error, len(m.modules))
	var wg sync.WaitGroup

	for _, mod := range m.modules {
		wg.Add(1)
		go func(mod MiniModule) {
			defer wg.Done()
			result, err := mod.Run(job)
			if err != nil {
				errChan <- fmt.Errorf("module %s failed: %w", mod.Name(), err)
				return
			}
			results[mod.Name()] = result
		}(mod)
	}

	// Wait for all modules to complete
	wg.Wait()
	close(errChan)

	// Check for errors
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	// Process aggregated results and update the database
	return m.saveResults(job.ScanID, results)
}

func (m *NetworkModule) saveResults(scanID uint, results map[string]interface{}) error {
	// Implementation of saving `results` to the database for `scanID`
	// This might involve calling the repository to save the combined result
	// for the scan
	fmt.Printf("Saving results for scan %s: %v\n", scanID, results)
	return nil
}

// Name returns the module name
func (m *NetworkModule) Name() string {
	return "Network"
}
