package modules

import (
	"crypto/x509"
	"fmt"
	"sync"

	"github.com/codevault-llc/humblebrag-api/internal/core/modules/network"
	"github.com/codevault-llc/humblebrag-api/internal/models/entities"
	"github.com/codevault-llc/humblebrag-api/internal/models/repository"
	"github.com/codevault-llc/humblebrag-api/internal/service"
	"github.com/codevault-llc/humblebrag-api/pkg/logger"
	"github.com/codevault-llc/humblebrag-api/pkg/utils"
	whoisparser "github.com/likexian/whois-parser"
	"go.uber.org/zap"
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
			&network.WhoisModule{},
			&network.DNSModule{},
			&network.CertificateModule{},
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
	networkModel := entities.NetworkModel{
		ScanID:      scanID,
		IPAddresses: results["IPLookup"].([]string),
		IPRanges:    results["IPRangeLookup"].([]string),
		HTTPHeaders: results["Headers"].([]string),
	}

	networkResponse, err := repository.
	if err != nil {
		logger.Log.Error("Failed to create network: %v", zap.Error(err))
		return err
	}

	whoisRecord := results["Whois"].(whoisparser.WhoisInfo)
	if whoisRecord.Registrar != nil {
		logger.Log.Info("Whois record:", zap.Any("whois", whoisRecord.Administrative))

		whois := entities.WhoisModel{
			NetworkId: networkResponse.ID,
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

		_, err = service.CreateWhois(whois)
		if err != nil {
			logger.Log.Error("Failed to create whois: %v", zap.Error(err))
			return err
		}
	}

	for _, certificate := range results["Certificates"].([]*x509.Certificate) {
		err := service.CreateCertificate(networkResponse.ID, *certificate)
		if err != nil {
			logger.Log.Error("Failed to create certificate: %v", zap.Error(err))
			return err
		}
	}

	dnsResults := results["DNS"].(network.DNSResults)
	dns := entities.DNSModel{
		NetworkId:   networkResponse.ID,
		CNAME:       dnsResults.CNAME,
		ARecords:    dnsResults.ARecords,
		AAAARecords: dnsResults.AAAARecords,
		MXRecords:   dnsResults.MXRecords,
		NSRecords:   dnsResults.NSRecords,
		TXTRecords:  dnsResults.TXTRecords,
		PTRRecord:   dnsResults.PTRRecord,
		DNSSEC:      dnsResults.DNSSEC,
	}

	_, err = service.CreateDNS(dns)
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
