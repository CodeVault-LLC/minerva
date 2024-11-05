package network

import (
	"crypto/x509"
	"fmt"
	"sync"

	generalEntities "github.com/codevault-llc/humblebrag-api/internal/core/models/entities"
	"github.com/codevault-llc/humblebrag-api/internal/network/models/entities"
	"github.com/codevault-llc/humblebrag-api/internal/network/models/repository"
	"github.com/codevault-llc/humblebrag-api/internal/network/modules"
	"github.com/codevault-llc/humblebrag-api/pkg/logger"
	"github.com/codevault-llc/humblebrag-api/pkg/types"
	"github.com/codevault-llc/humblebrag-api/pkg/utils"
	whoisparser "github.com/likexian/whois-parser"
	"go.uber.org/zap"
)

type MiniModule interface {
	Run(job generalEntities.JobModel) (interface{}, error) // Executes the mini-module logic
	Name() string                                          // Returns the mini-module name
}

// NetworkModule orchestrates network-related scans through sub-modules
type NetworkModule struct {
	modules []MiniModule
}

func NewNetworkModule() *NetworkModule {
	return &NetworkModule{
		modules: []MiniModule{
			&modules.IPLookupModule{},
			&modules.IPRangeLookupModule{},
			&modules.HeaderModule{},
			&modules.WhoisModule{},
			&modules.DNSModule{},
			&modules.CertificateModule{},
		},
	}
}

// Execute runs the Network-specific scan logic
func (m *NetworkModule) Execute(job generalEntities.JobModel, website types.WebsiteAnalysis) error {
	var wg sync.WaitGroup
	var mu sync.Mutex
	results := make(map[string]interface{})
	errChan := make(chan error, len(m.modules))

	for _, mod := range m.modules {
		wg.Add(1)
		logger.Log.Info("Running module", zap.String("module", mod.Name()))
		go func(mod MiniModule) {
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
		ScanId:      scanID,
		IpAddresses: results["IPLookup"].([]string),
		IpRanges:    results["IPRangeLookup"].([]string),
		HttpHeaders: results["Header"].([]string),
	}

	networkResponse, err := repository.NetworkRepository.Create(networkModel)
	if err != nil {
		logger.Log.Error("Failed to create network: %v", zap.Error(err))
		return err
	}

	whoisRecord := results["Whois"].(whoisparser.WhoisInfo)
	if whoisRecord.Registrar != nil {
		logger.Log.Info("Whois record:", zap.Any("whois", whoisRecord.Administrative))

		whois := entities.WhoisModel{
			NetworkId: networkResponse.Id,
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

		err := repository.WhoisRepository.SaveWhoisResult(whois)
		if err != nil {
			logger.Log.Error("Failed to create whois: %v", zap.Error(err))
			return err
		}
	}

	for _, certificate := range results["Certificate"].([]*x509.Certificate) {
		_, err := repository.CertificateRepository.Create(networkResponse.Id, *certificate)
		if err != nil {
			logger.Log.Error("Failed to create certificate: %v", zap.Error(err))
			return err
		}
	}

	dnsResults := results["DNS"].(modules.DNSResults)
	dns := entities.DnsModel{
		NetworkId:   networkResponse.Id,
		Cname:       dnsResults.CNAME,
		ARecords:    dnsResults.ARecords,
		AAAARecords: dnsResults.AAAARecords,
		MxRecords:   dnsResults.MXRecords,
		NsRecords:   dnsResults.NSRecords,
		TxtRecords:  dnsResults.TXTRecords,
		PtrRecord:   dnsResults.PTRRecord,
		Dnssec:      dnsResults.DNSSEC,
	}

	err = repository.DnsRepository.SaveDnsResult(dns)
	if err != nil {
		logger.Log.Error("Failed to create DNS: %v", zap.Error(err))
		return err
	}
	return nil
}

// Name returns the module name
func (m *NetworkModule) Name() string {
	return "Network"
}
