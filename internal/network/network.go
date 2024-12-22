package network

import (
	"crypto/x509"
	"fmt"
	"sync"

	"github.com/codevault-llc/minerva/internal/common"
	generalEntities "github.com/codevault-llc/minerva/internal/core/models/entities"
	"github.com/codevault-llc/minerva/internal/database"
	"github.com/codevault-llc/minerva/internal/network/models/entities"
	"github.com/codevault-llc/minerva/internal/network/models/repository"
	"github.com/codevault-llc/minerva/internal/network/modules"
	"github.com/codevault-llc/minerva/pkg/logger"
	"github.com/codevault-llc/minerva/pkg/utils"
	"github.com/google/uuid"
	whoisparser "github.com/likexian/whois-parser"
	"go.uber.org/zap"
)

// NetworkModule orchestrates network-related scans through sub-modules
type NetworkModule struct {
	modules         []common.MiniModule
	runtimeLocation common.RuntimeLocation
	repository      *repository.NetworkRepo
}

func NewNetworkModule(runtimeLocation common.RuntimeLocation, db *database.Database) *NetworkModule {
	repository.NetworkRepository = repository.NewNetworkRepository(db)

	return &NetworkModule{
		modules: []common.MiniModule{
			&modules.IPLookupModule{},
			&modules.IPRangeLookupModule{},
			&modules.HeaderModule{},
			&modules.WhoisModule{},
			&modules.DNSModule{},
			&modules.CertificateModule{},
		},
		runtimeLocation: runtimeLocation,
		repository:      repository.NetworkRepository,
	}
}

// Execute runs the Network-specific scan logic
func (m *NetworkModule) Execute(job generalEntities.JobModel, website *common.WebsiteAnalysis) error {
	var wg sync.WaitGroup
	var mu sync.Mutex
	results := make(map[string]interface{})
	errChan := make(chan error, len(m.modules))

	for _, mod := range m.modules {
		wg.Add(1)
		go func(mod common.MiniModule) {
			defer wg.Done()
			result, err := mod.Run(job)
			if err != nil {
				errChan <- fmt.Errorf("module %s failed: %w", mod.Name(), err)
			}
			mu.Lock()
			results[mod.Name()] = result
			mu.Unlock()
		}(mod)
	}

	// Wait for all modules to complete
	wg.Wait()
	close(errChan)

	// Process aggregated results and update the database
	return m.saveResults(job.ScanID, results)
}

func (m *NetworkModule) saveResults(scanID string, results map[string]interface{}) error {
	dnsResults := results["DNS"].(modules.DNSResults)
	dnsModel := entities.DnsModel{
		Cname:       dnsResults.CNAME,
		ARecords:    dnsResults.ARecords,
		AAAARecords: dnsResults.AAAARecords,
		MxRecords:   dnsResults.MXRecords,
		NsRecords:   dnsResults.NSRecords,
		TxtRecords:  dnsResults.TXTRecords,
		PtrRecord:   dnsResults.PTRRecord,
		Dnssec:      dnsResults.DNSSEC,
	}

	whoisRecord := results["Whois"].(whoisparser.WhoisInfo)
	certificates := results["Certificate"].([]*x509.Certificate)

	networkModel := entities.NetworkModel{
		Id:          uuid.New().String(),
		IpAddresses: results["IPLookup"].([]string),
		IpRanges:    results["IPRangeLookup"].([]string),
		HttpHeaders: results["Header"].([]string),

		WhoisModel:       whoisRecord,
		DnsModel:         dnsModel,
		CertificateModel: certificates,

		CreatedAt: utils.GetCurrentTime(),
		UpdatedAt: utils.GetCurrentTime(),
	}

	err := m.repository.Create(scanID, networkModel)
	if err != nil {
		logger.Log.Error("Failed to create network: %v", zap.Error(err))
		return err
	}

	return nil
}

// Name returns the module name
func (m *NetworkModule) Name() string {
	return "Network"
}

func (m *NetworkModule) RuntimeLocation() common.RuntimeLocation {
	return m.runtimeLocation
}
