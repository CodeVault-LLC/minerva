package network

import (
	"crypto/x509"
	"fmt"
	"sync"

	"github.com/codevault-llc/humblebrag-api/internal/database/models"
	"github.com/codevault-llc/humblebrag-api/internal/service"
	"github.com/codevault-llc/humblebrag-api/pkg/logger"
	"github.com/codevault-llc/humblebrag-api/pkg/utils"
	whoisparser "github.com/likexian/whois-parser"
	"go.uber.org/zap"
)

func NetworkModule(scanId uint, url string) {
	ipAddrChan := make(chan []string)
	ipRangeChan := make(chan []string)
	dnsResultChan := make(chan DNSResults)
	whoisChan := make(chan whoisparser.WhoisInfo)
	certificateChan := make(chan []*x509.Certificate)
	headerChan := make(chan []string)

	var wg sync.WaitGroup
	wg.Add(6)

	go func() {
		defer wg.Done()
		ipAddresses, err := ScanIP(url)
		if err != nil {
			ipAddrChan <- nil // Handle error appropriately
			logger.Log.Error("Failed to scan IP: %v", zap.Error(err))
		} else {
			ipAddrChan <- ipAddresses
		}
	}()

	go func() {
		defer wg.Done()
		ipRanges, err := ScanIPRange(url)
		if err != nil {
			ipRangeChan <- nil // Handle error appropriately
			logger.Log.Error("Failed to scan IP range: %v", zap.Error(err))
		} else {
			ipRangeChan <- ipRanges
		}
	}()

	go func() {
		defer wg.Done()
		dnsResults, err := GetDNSScan(url)
		if err != nil {
			dnsResultChan <- DNSResults{} // Handle error appropriately
			logger.Log.Error("Failed to scan DNS: %v", zap.Error(err))
		} else {
			dnsResultChan <- dnsResults
		}
	}()

	go func() {
		defer wg.Done()
		whoisRecord, err := ScanWhois(utils.ConvertURLToDomain(url))
		if err != nil || whoisRecord.Domain.Name == "" {
			// Handle the case where the domain doesn't exist or the whois data is incomplete
			whoisChan <- whoisparser.WhoisInfo{}
			logger.Log.Error("Failed to scan whois: %v", zap.Error(err))
		} else {
			whoisChan <- whoisRecord
		}
	}()

	go func() {
		defer wg.Done()
		certifiate, err := GetCertificateWebsite(url, 443)
		if err != nil {
			certificateChan <- nil // Handle error appropriately
			logger.Log.Error("Failed to get certificate: %v", zap.Error(err))
		} else {
			certificateChan <- certifiate
		}
	}()

	go func() {
		defer wg.Done()
		header, err := getHeaders(url)
		if err != nil {
			headerChan <- nil // Handle error appropriately
			logger.Log.Error("Failed to get headers: %v", zap.Error(err))
		} else {
			httpHeaders := make([]string, 0)
			for key, value := range header.Headers {
				httpHeaders = append(httpHeaders, fmt.Sprintf("%s: %s", key, value))
			}
			headerChan <- httpHeaders
		}
	}()

	go func() {
		wg.Wait()
		close(ipAddrChan)
		close(ipRangeChan)
		close(dnsResultChan)
		close(whoisChan)
		close(certificateChan)
		close(headerChan)
	}()

	ipAddresses := <-ipAddrChan
	ipRanges := <-ipRangeChan
	dnsResults := <-dnsResultChan
	whoisRecord := <-whoisChan
	certifiate := <-certificateChan
	headers := <-headerChan

	network := models.NetworkModel{
		ScanID:      scanId,
		IPAddresses: ipAddresses,
		IPRanges:    ipRanges,
		HTTPHeaders: headers,
	}

	networkResponse, err := service.CreateNetwork(network)
	if err != nil {
		logger.Log.Error("Failed to create network: %v", zap.Error(err))
		return
	}

	if whoisRecord.Registrar != nil {
		whois := models.WhoisModel{
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

			AdminName:       utils.SafeString(whoisRecord.Administrative.Name),
			AdminEmail:      utils.SafeString(whoisRecord.Administrative.Email),
			AdminPhone:      utils.SafeString(whoisRecord.Administrative.Phone),
			AdminOrg:        utils.SafeString(whoisRecord.Administrative.Organization),
			AdminCity:       utils.SafeString(whoisRecord.Administrative.City),
			AdminPostalCode: utils.SafeString(whoisRecord.Administrative.PostalCode),
			AdminCountry:    utils.SafeString(whoisRecord.Administrative.Country),

			Updated: whoisRecord.Domain.UpdatedDate,
			Created: whoisRecord.Domain.CreatedDate,
			Expires: whoisRecord.Domain.ExpirationDate,
		}

		_, err = service.CreateWhois(whois)
		if err != nil {
			logger.Log.Error("Failed to create whois: %v", zap.Error(err))
		}
	}

	for _, certificate := range certifiate {
		err := service.CreateCertificate(networkResponse.ID, *certificate)
		if err != nil {
			logger.Log.Error("Failed to create certificate: %v", zap.Error(err))
			return
		}
	}

	dns := models.DNSModel{
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
	}
}
