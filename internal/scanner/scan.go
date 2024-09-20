package scanner

import (
	"crypto/x509"
	"fmt"
	"net/http"
	"sync"

	"github.com/codevault-llc/humblebrag-api/internal/scanner/certificate"
	"github.com/codevault-llc/humblebrag-api/internal/scanner/http_req"
	"github.com/codevault-llc/humblebrag-api/internal/scanner/network"
	"github.com/codevault-llc/humblebrag-api/internal/scanner/secrets"
	"github.com/codevault-llc/humblebrag-api/internal/scanner/security"
	"github.com/codevault-llc/humblebrag-api/internal/scanner/websites"
	"github.com/codevault-llc/humblebrag-api/internal/service"
	"github.com/codevault-llc/humblebrag-api/internal/updater"
	"github.com/codevault-llc/humblebrag-api/models"
	"github.com/codevault-llc/humblebrag-api/pkg/logger"
	"github.com/codevault-llc/humblebrag-api/pkg/parsers"
	"github.com/codevault-llc/humblebrag-api/pkg/utils"
	"github.com/codevault-llc/humblebrag-api/types"
	whoisparser "github.com/likexian/whois-parser"
)

type WebsiteScan struct {
	Website      models.ScanResponse
	IPAddresses  []string
	IPRanges     []string
	HTTPHeaders  []string
	Certificates []*x509.Certificate
	Secrets      []utils.RegexReturn
	GetDNSScan   network.DNSResults
	FoundLists   []types.List
	WhoisRecord  whoisparser.WhoisInfo
}

type NetworkScan struct {
	IPAddresses []string
	IPRanges    []string
	GetDNSScan  network.DNSResults
	WhoisRecord whoisparser.WhoisInfo
}

type SecurityScan struct {
	corsIssues      []string
	headerScan      security.ScanSecurityHead
	protocolSupport []string
}

func ScanWebsite(url string, userId uint) (models.ScanModel, error) {
	// Start timing the function
	timer := logger.StartTimer()

	var wg sync.WaitGroup
	_ = make(chan models.ScanResponse, 1) // Buffer of 1 to avoid blocking
	scanUpdateChan := make(chan error)

	// Fetch website content first to immediately return a basic scan response
	logger.Log.Info("Starting website scan for URL: %s", url)
	website, err := websites.ScanWebsite(url)
	if err != nil {
		logger.Log.Error("Failed to scan website: %v", err)
		return models.ScanModel{}, err
	}

	// Create initial scan object
	websiteScan := WebsiteScan{
		Website: website,
		// Basic initial fields like the main website content can be filled here
	}

	logger.Log.Info("Saving initial scan result for URL: %s", url)
	// Save the initial scan result to DB
	scan, err := saveScan(websiteScan, userId)
	if err != nil {
		logger.Log.Error("Failed to save initial scan result: %v", err)
		return models.ScanModel{}, err
	}

	// Start other scans asynchronously, and update scan records as they finish
	wg.Add(1)
	go func() {
		defer wg.Done()
		// Fetch HTTP response asynchronously
		logger.Log.Info("Fetching HTTP response for URL: %s", url)
		httpResponse, err := http_req.GetHTTPResponse(url)
		if err != nil {
			logger.Log.Error("Failed to fetch HTTP response: %v", err)
		} else {
			logger.Log.Info("HTTP response fetched for URL: %s", url)
			httpHeaders := make([]string, 0)
			for key, value := range httpResponse.Headers {
				httpHeaders = append(httpHeaders, fmt.Sprintf("%s: %s", key, value))
			}
			websiteScan.HTTPHeaders = httpHeaders

			// Update scan with HTTP headers
			_, err = updateScan(websiteScan, scan.ID) // Asynchronous update
			if err != nil {
				logger.Log.Error("Failed to update scan with found lists: %v", err)
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		// Fetch certificates asynchronously
		logger.Log.Info("Fetching certificates for URL: %s", url)
		certificates, err := certificate.GetCertificateWebsite(url, 443)
		if err != nil {
			logger.Log.Error("Failed to fetch certificates: %v", err)
		} else {
			logger.Log.Info("Certificates fetched for URL: %s", url)
			websiteScan.Certificates = certificates

			// Update scan with certificates
			_, err = updateScan(websiteScan, scan.ID) // Asynchronous update
			if err != nil {
				logger.Log.Error("Failed to update scan with found lists: %v", err)
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		// Scan secrets asynchronously (dependent on website scan)
		logger.Log.Info("Scanning secrets for URL: %s", url)
		secretsFound := secrets.ScanSecrets(website.Scripts)
		websiteScan.Secrets = secretsFound
		logger.Log.Info("Secrets scan completed for URL: %s", url)

		// Update scan with secrets found
		_, err = updateScan(websiteScan, scan.ID) // Asynchronous update
		if err != nil {
			logger.Log.Error("Failed to update scan with found lists: %v", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		// Perform network scan asynchronously
		logger.Log.Info("Performing network scan for URL: %s", url)
		networkScan := scanNetwork(url)
		websiteScan.IPAddresses = networkScan.IPAddresses
		websiteScan.IPRanges = networkScan.IPRanges
		websiteScan.GetDNSScan = networkScan.GetDNSScan
		websiteScan.WhoisRecord = networkScan.WhoisRecord
		logger.Log.Info("Network scan completed for URL: %s", url)

		// Update scan with network details
		_, err = updateScan(websiteScan, scan.ID) // Asynchronous update
		if err != nil {
			logger.Log.Error("Failed to update scan with found lists: %v", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		// Compare found lists asynchronously
		logger.Log.Info("Comparing found lists for URL: %s", url)
		foundLists := updater.CompareValues(utils.ConvertURLToDomain(url), parsers.Domain)
		websiteScan.FoundLists = foundLists
		logger.Log.Info("List comparison completed for URL: %s", url)

		// Update scan with found lists
		_, err = updateScan(websiteScan, scan.ID) // Asynchronous update
		if err != nil {
			logger.Log.Error("Failed to update scan with found lists: %v", err)
		}
	}()

	// Optionally wait for all updates to finish (non-blocking version)
	go func() {
		wg.Wait()
		close(scanUpdateChan)
	}()

	// Stop the timer and log the time taken
	timer.Stop(logger.Log, "Initial ScanWebsite execution")

	// Return the initial scan immediately (with basic website data)
	return scan, nil
}

func scanSecurity(addr string, headers http.Header) SecurityScan {
	corsIssues := security.ScanCors(headers)
	headerScan := security.ScanSecurityHeaders(headers)
	protocolSupport, err := security.ScanProtocolSupport(addr)

	fmt.Println("Protocol Support", protocolSupport)

	if err != nil {
		fmt.Println("Failed to scan protocol support", err)
	}

	return SecurityScan{
		corsIssues:      corsIssues,
		headerScan:      headerScan,
		protocolSupport: protocolSupport,
	}
}

func scanNetwork(url string) NetworkScan {
	ipAddresses, _ := network.ScanIP(url)
	ipRanges, _ := network.ScanIPRange(url)
	dnsResults, _ := network.GetDNSScan(url)
	whoisRecord, _ := network.ScanWhois(utils.ConvertURLToDomain(url))

	return NetworkScan{
		IPAddresses: ipAddresses,
		IPRanges:    ipRanges,
		GetDNSScan:  dnsResults,
		WhoisRecord: whoisRecord,
	}
}

func saveScan(scan WebsiteScan, userId uint) (models.ScanModel, error) {
	scanModel := models.ScanModel{
		WebsiteUrl:  scan.Website.WebsiteUrl,
		WebsiteName: scan.Website.WebsiteName,

		UserID: userId,

		Sha256: utils.SHA256(scan.Website.WebsiteUrl),
		SHA1:   utils.SHA1(scan.Website.WebsiteUrl),
		MD5:    utils.MD5(scan.Website.WebsiteUrl),

		Status: models.ScanStatusComplete,
	}

	// Create Scan
	scanResponse, err := service.CreateScan(scanModel)
	if err != nil {
		fmt.Println("Failed to create scan", err)
		return models.ScanModel{}, err
	}

	// Create Findings
	service.CreateFindings(scanResponse.ID, scan.Secrets)

	// Create Contents
	for _, script := range scan.Website.Scripts {
		content := models.ContentModel{
			ScanID:  scanResponse.ID,
			Name:    script.Src,
			Content: script.Content,
		}

		_, err := service.CreateContent(content)
		if err != nil {
			return models.ScanModel{}, err
		}
	}

	// Create Networks
	network := models.NetworkModel{
		ScanID:       scanResponse.ID,
		IPAddresses:  scan.IPAddresses,
		HTTPHeaders:  scan.HTTPHeaders,
		IPRanges:     scan.IPRanges,
		DNSNames:     scan.GetDNSScan.CNAME,
		PermittedDNS: scan.GetDNSScan.Permitted,
		ExcludedDNS:  scan.GetDNSScan.Excluded,
	}

	networkResponse, err := service.CreateNetwork(network)
	if err != nil {
		return models.ScanModel{}, err
	}

	// Create Certificates
	for _, certificate := range scan.Certificates {
		err := service.CreateCertificate(networkResponse.ID, *certificate)
		if err != nil {
			fmt.Println("Failed to create certificate", err)
			return models.ScanModel{}, err
		}
	}

	// Create Whois
	whois := models.WhoisModel{
		NetworkId: networkResponse.ID,
		Status: func() string {
			if len(scan.WhoisRecord.Domain.Status) > 0 {
				return scan.WhoisRecord.Domain.Status[0]
			}
			return ""
		}(),

		DomainName:  scan.WhoisRecord.Domain.Name,
		Registrar:   scan.WhoisRecord.Registrar.Name,
		Email:       scan.WhoisRecord.Registrant.Email,
		Phone:       scan.WhoisRecord.Registrant.Phone,
		NameServers: scan.WhoisRecord.Domain.NameServers,

		RegistrantName:       scan.WhoisRecord.Registrant.Name,
		RegistrantCity:       scan.WhoisRecord.Registrant.City,
		RegistrantPostalCode: scan.WhoisRecord.Registrant.PostalCode,
		RegistrantCountry:    scan.WhoisRecord.Registrant.Country,
		RegistrantEmail:      scan.WhoisRecord.Registrant.Email,
		RegistrantPhone:      scan.WhoisRecord.Registrant.Phone,
		RegistrantOrg:        scan.WhoisRecord.Registrant.Organization,

		AdminName:       scan.WhoisRecord.Administrative.Name,
		AdminEmail:      scan.WhoisRecord.Administrative.Email,
		AdminPhone:      scan.WhoisRecord.Administrative.Phone,
		AdminOrg:        scan.WhoisRecord.Administrative.Organization,
		AdminCity:       scan.WhoisRecord.Administrative.City,
		AdminPostalCode: scan.WhoisRecord.Administrative.PostalCode,
		AdminCountry:    scan.WhoisRecord.Administrative.Country,

		Updated: scan.WhoisRecord.Domain.UpdatedDate,
		Created: scan.WhoisRecord.Domain.CreatedDate,
		Expires: scan.WhoisRecord.Domain.ExpirationDate,
	}

	_, err = service.CreateWhois(whois)
	if err != nil {
		return models.ScanModel{}, err
	}

	for _, list := range scan.FoundLists {
		listModel := models.ListModel{
			ScanID: scanResponse.ID,
			ListID: list.ListID,
		}

		_, err := service.CreateList(listModel)
		if err != nil {
			return models.ScanModel{}, err
		}
	}

	return scanResponse, nil
}

func updateScan(scan WebsiteScan, scanId uint) (models.ScanModel, error) {
	scanModel := models.ScanModel{
		WebsiteUrl:  scan.Website.WebsiteUrl,
		WebsiteName: scan.Website.WebsiteName,

		Sha256: utils.SHA256(scan.Website.WebsiteUrl),
		SHA1:   utils.SHA1(scan.Website.WebsiteUrl),
		MD5:    utils.MD5(scan.Website.WebsiteUrl),

		Status: models.ScanStatusComplete,
	}
	scanModel.ID = scanId

	// Update Scan
	scanResponse, err := service.UpdateScan(scanModel)
	if err != nil {
		return models.ScanModel{}, err
	}

	// Update Findings
	service.UpdateFindings(scanResponse.ID, scan.Secrets)

	// Update Contents
	service.DeleteContents(scanResponse.ID)
	for _, script := range scan.Website.Scripts {
		content := models.ContentModel{
			ScanID:  scanResponse.ID,
			Name:    script.Src,
			Content: script.Content,
		}

		_, err := service.CreateContent(content)
		if err != nil {
			return models.ScanModel{}, err
		}
	}

	// Update Networks
	network := models.NetworkModel{
		ScanID:       scanResponse.ID,
		IPAddresses:  scan.IPAddresses,
		HTTPHeaders:  scan.HTTPHeaders,
		IPRanges:     scan.IPRanges,
		DNSNames:     scan.GetDNSScan.CNAME,
		PermittedDNS: scan.GetDNSScan.Permitted,
		ExcludedDNS:  scan.GetDNSScan.Excluded,
	}

	networkResponse, err := service.UpdateNetwork(network)
	if err != nil {
		return models.ScanModel{}, err
	}

	// Update Certificates
	service.DeleteCertificates(networkResponse.ID)
	for _, certificate := range scan.Certificates {
		err := service.CreateCertificate(networkResponse.ID, *certificate)
		if err != nil {
			return models.ScanModel{}, err
		}
	}

	// Create Whois
	whois := models.WhoisModel{
		NetworkId: networkResponse.ID,
		Status: func() string {
			if len(scan.WhoisRecord.Domain.Status) > 0 {
				return scan.WhoisRecord.Domain.Status[0]
			}
			return ""
		}(),

		DomainName:  scan.WhoisRecord.Domain.Name,
		Registrar:   scan.WhoisRecord.Registrar.Name,
		Email:       scan.WhoisRecord.Registrant.Email,
		Phone:       scan.WhoisRecord.Registrant.Phone,
		NameServers: scan.WhoisRecord.Domain.NameServers,

		RegistrantName:       scan.WhoisRecord.Registrant.Name,
		RegistrantCity:       scan.WhoisRecord.Registrant.City,
		RegistrantPostalCode: scan.WhoisRecord.Registrant.PostalCode,
		RegistrantCountry:    scan.WhoisRecord.Registrant.Country,
		RegistrantEmail:      scan.WhoisRecord.Registrant.Email,
		RegistrantPhone:      scan.WhoisRecord.Registrant.Phone,
		RegistrantOrg:        scan.WhoisRecord.Registrant.Organization,

		AdminName:       scan.WhoisRecord.Administrative.Name,
		AdminEmail:      scan.WhoisRecord.Administrative.Email,
		AdminPhone:      scan.WhoisRecord.Administrative.Phone,
		AdminOrg:        scan.WhoisRecord.Administrative.Organization,
		AdminCity:       scan.WhoisRecord.Administrative.City,
		AdminPostalCode: scan.WhoisRecord.Administrative.PostalCode,
		AdminCountry:    scan.WhoisRecord.Administrative.Country,

		Updated: scan.WhoisRecord.Domain.UpdatedDate,
		Created: scan.WhoisRecord.Domain.CreatedDate,
		Expires: scan.WhoisRecord.Domain.ExpirationDate,
	}

	_, err = service.CreateWhois(whois)
	if err != nil {
		return models.ScanModel{}, err
	}

	for _, list := range scan.FoundLists {
		listModel := models.ListModel{
			ScanID: scanResponse.ID,
			ListID: list.ListID,
		}

		_, err := service.CreateList(listModel)
		if err != nil {
			return models.ScanModel{}, err
		}
	}

	return scanResponse, nil
}
