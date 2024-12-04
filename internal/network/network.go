package network

import (
	"crypto/x509"
	"fmt"
	"sync"

	"github.com/codevault-llc/minerva/internal/common"
	generalEntities "github.com/codevault-llc/minerva/internal/core/models/entities"
	"github.com/codevault-llc/minerva/internal/network/models/entities"
	"github.com/codevault-llc/minerva/internal/network/models/repository"
	"github.com/codevault-llc/minerva/internal/network/modules"
	"github.com/codevault-llc/minerva/pkg/logger"
	"github.com/codevault-llc/minerva/pkg/utils"
	whoisparser "github.com/likexian/whois-parser"
	"go.uber.org/zap"
)

// NetworkModule orchestrates network-related scans through sub-modules
type NetworkModule struct {
	modules         []common.MiniModule
	runtimeLocation common.RuntimeLocation
}

func NewNetworkModule(runtimeLocation common.RuntimeLocation) *NetworkModule {
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
		logger.Log.Info("Running module", zap.String("module", mod.Name()))
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

func (m *NetworkModule) saveResults(scanID uint, results map[string]interface{}) error {
	networkModel := entities.NetworkModel{
		IpAddresses: results["IPLookup"].([]string),
		IpRanges:    results["IPRangeLookup"].([]string),
		HttpHeaders: results["Header"].([]string),
	}

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
	if whoisRecord.Registrar != nil {
		logger.Log.Info("Whois record:", zap.Any("whois", whoisRecord.Administrative))
	}

	whoisModel := entities.WhoisModel{
		Status: func() string {
			if len(whoisRecord.Domain.Status) > 0 {
				return whoisRecord.Domain.Status[0]
			}
			return ""
		}(),
		DomainName:  whoisRecord.Domain.Name,
		Registrar:   utils.SafeString(whoisRecord.Registrar.Name),
		Email:       utils.SafeString(whoisRecord.Registrant.Email),
		Phone:       utils.SafeString(whoisRecord.Registrant.Phone),
		NameServers: whoisRecord.Domain.NameServers,

		RegistrantName:       utils.SafeString(whoisRecord.Registrant.Name),
		RegistrantCity:       utils.SafeString(whoisRecord.Registrant.City),
		RegistrantPostalCode: utils.SafeString(whoisRecord.Registrant.PostalCode),
		RegistrantCountry:    utils.SafeString(whoisRecord.Registrant.Country),
		RegistrantEmail:      utils.SafeString(whoisRecord.Registrant.Email),
		RegistrantPhone:      utils.SafeString(whoisRecord.Registrant.Phone),
		RegistrantOrg:        utils.SafeString(whoisRecord.Registrant.Organization),
		AdminName: func() string {
			if whoisRecord.Administrative != nil && whoisRecord.Administrative.Name != "" {
				return whoisRecord.Administrative.Name
			}
			return ""
		}(),
		AdminEmail: func() string {
			if whoisRecord.Administrative != nil {
				return utils.SafeString(whoisRecord.Administrative.Email)
			}
			return ""
		}(),
		AdminPhone: func() string {
			if whoisRecord.Administrative != nil {
				return utils.SafeString(whoisRecord.Administrative.Phone)
			}
			return ""
		}(),
		AdminOrg: func() string {
			if whoisRecord.Administrative != nil {
				return utils.SafeString(whoisRecord.Administrative.Organization)
			}
			return ""
		}(),
		AdminCity: func() string {
			if whoisRecord.Administrative != nil {
				return utils.SafeString(whoisRecord.Administrative.City)
			}
			return ""
		}(),
		AdminPostalCode: func() string {
			if whoisRecord.Administrative != nil {
				return utils.SafeString(whoisRecord.Administrative.PostalCode)
			}
			return ""
		}(),
		AdminCountry: func() string {
			if whoisRecord.Administrative != nil {
				return utils.SafeString(whoisRecord.Administrative.Country)
			}
			return ""
		}(),

		Updated: whoisRecord.Domain.UpdatedDate,
		Created: whoisRecord.Domain.CreatedDate,
		Expires: whoisRecord.Domain.ExpirationDate,
	}

	certificates := results["Certificate"].([]x509.Certificate)

	networkId, err := repository.NetworkRepository.Create(scanID, networkModel, whoisModel, dnsModel, certificates)
	if err != nil {
		logger.Log.Error("Failed to create network: %v", zap.Error(err))
		return err
	}

	logger.Log.Info("Network created", zap.Uint("networkId", networkId))

	return nil
}

// Name returns the module name
func (m *NetworkModule) Name() string {
	return "Network"
}

func (m *NetworkModule) RuntimeLocation() common.RuntimeLocation {
	return m.runtimeLocation
}
