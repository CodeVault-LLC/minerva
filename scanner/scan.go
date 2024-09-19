package scanner

import (
	"crypto/x509"
	"fmt"
	"net/http"

	"github.com/codevault-llc/humblebrag-api/models"
	"github.com/codevault-llc/humblebrag-api/parsers"
	"github.com/codevault-llc/humblebrag-api/scanner/certificate"
	"github.com/codevault-llc/humblebrag-api/scanner/http_req"
	"github.com/codevault-llc/humblebrag-api/scanner/network"
	"github.com/codevault-llc/humblebrag-api/scanner/secrets"
	"github.com/codevault-llc/humblebrag-api/scanner/security"
	"github.com/codevault-llc/humblebrag-api/scanner/websites"
	"github.com/codevault-llc/humblebrag-api/service"
	"github.com/codevault-llc/humblebrag-api/types"
	"github.com/codevault-llc/humblebrag-api/updater"
	"github.com/codevault-llc/humblebrag-api/utils"
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
	website, _ := websites.ScanWebsite(url)
	httpResponse, _ := http_req.GetHTTPResponse(url)
	certificates, _ := certificate.GetCertificateWebsite(url, 443)
	secretsFound := secrets.ScanSecrets(website.Scripts)

	addr := fmt.Sprintf("%s:%d", utils.ConvertURLToDomain(url), 443)
	go scanSecurity(addr, httpResponse.Headers)

	foundLists := updater.CompareValues(utils.ConvertURLToDomain(url), parsers.Domain)
	networkScan := scanNetwork(url)

	// Create HTTP Headers
	httpHeaders := make([]string, 0)
	for key, value := range httpResponse.Headers {
		httpHeaders = append(httpHeaders, fmt.Sprintf("%s: %s", key, value))
	}

	websiteScan := WebsiteScan{
		Website:      website,
		IPAddresses:  networkScan.IPAddresses,
		IPRanges:     networkScan.IPRanges,
		HTTPHeaders:  httpHeaders,
		Certificates: certificates,
		Secrets:      secretsFound,
		GetDNSScan:   networkScan.GetDNSScan,
		FoundLists:   foundLists,
		WhoisRecord:  networkScan.WhoisRecord,
	}

	scan, err := saveScan(websiteScan, userId)
	if err != nil {
		return models.ScanModel{}, err
	}

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

		Sha256: fmt.Sprintf("%x", utils.SHA256(scan.Website.WebsiteUrl)),
		SHA1:   fmt.Sprintf("%x", utils.SHA1(scan.Website.WebsiteUrl)),
		MD5:    fmt.Sprintf("%x", utils.MD5(scan.Website.WebsiteUrl)),

		Status: models.ScanStatusComplete,
	}

	// Create Scan
	scanResponse, err := service.CreateScan(scanModel)
	if err != nil {
		fmt.Println("Failed to create scan", err)
		return models.ScanModel{}, err
	}

	// Create Certificates
	for _, certificate := range scan.Certificates {
		err := service.CreateCertificate(scanResponse.ID, *certificate)
		if err != nil {
			fmt.Println("Failed to create certificate", err)
			return models.ScanModel{}, err
		}
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

	// Create Whois
	whois := models.WhoisModel{
		ScanID: scanResponse.ID,
		Status: scan.WhoisRecord.Domain.Status[0],

		DomainName: scan.WhoisRecord.Domain.Name,
		Registrar:  scan.WhoisRecord.Registrar.Name,
		Email:      scan.WhoisRecord.Registrant.Email,
		Phone:      scan.WhoisRecord.Registrant.Phone,
		NameServer: scan.WhoisRecord.Domain.NameServers,

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

	_, err = service.CreateNetwork(network)
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
