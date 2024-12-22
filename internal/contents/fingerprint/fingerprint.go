package fingerprint

import (
	"sync"

	"github.com/codevault-llc/minerva/config"
	"github.com/codevault-llc/minerva/pkg/types"
	"github.com/codevault-llc/minerva/pkg/utils"
)

type FingerprintModule struct {
}

func NewFingerprintModule() *FingerprintModule {
	return &FingerprintModule{}
}

func (m *FingerprintModule) Execute(filename string) []utils.RegexReturn {
	var results []utils.RegexReturn

	var wg sync.WaitGroup
	var mu sync.Mutex

	concurrencyLimit := make(chan struct{}, 10)

	for _, rule := range config.Config.Fingerprints {
		concurrencyLimit <- struct{}{}
		wg.Add(1)

		go func(rule types.Fingerprint) {
			defer wg.Done()
			defer func() { <-concurrencyLimit }()

			var scriptResults []utils.Match
			matches := utils.GenericScanFingerprint(rule, filename)
			if len(matches) > 0 {
				scriptResults = append(scriptResults, matches...)
			}

			if len(scriptResults) > 0 {
				mu.Lock()
				results = append(results, utils.RegexReturn{Name: rule.Name, Matches: scriptResults, Description: rule.Description})
				mu.Unlock()
			}
		}(rule)
	}

	wg.Wait()
	return results
}
