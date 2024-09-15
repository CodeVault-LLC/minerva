package scanner

import (
	"crypto/x509"
	"fmt"

	"github.com/codevault-llc/humblebrag-api/controller"
	"github.com/codevault-llc/humblebrag-api/models"
	"github.com/codevault-llc/humblebrag-api/scanner/certificate"
	"github.com/codevault-llc/humblebrag-api/scanner/http_req"
	"github.com/codevault-llc/humblebrag-api/scanner/ip"
	"github.com/codevault-llc/humblebrag-api/scanner/secrets"
	"github.com/codevault-llc/humblebrag-api/scanner/websites"
	"github.com/codevault-llc/humblebrag-api/utils"
)

type WebsiteScan struct {
	Website      models.ScanResponse
	IPAddresses  []string
	HTTPHeaders  map[string][]string
	Certificates []*x509.Certificate
	Secrets      []utils.RegexReturn
}

func ScanWebsite(url string) (models.Scan, error) {
	website, _ := websites.ScanWebsite(url)
	ipAddresses, _ := ip.ScanIP(url)
	httpHeaders, _ := http_req.ScanHTTPHeaders(url)
	certificates, _ := certificate.GetCertificateWebsite(url, 443)
	secretsFound := secrets.ScanSecrets(website.Scripts)

	websiteScan := WebsiteScan{
		Website:      website,
		IPAddresses:  ipAddresses,
		HTTPHeaders:  httpHeaders,
		Certificates: certificates,
		Secrets:      secretsFound,
	}

	scan, err := saveScan(websiteScan)
	if err != nil {
		return models.Scan{}, err
	}

	return scan, nil
}

func saveScan(scan WebsiteScan) (models.Scan, error) {
	scanModel := models.Scan{
		WebsiteUrl:  scan.Website.WebsiteUrl,
		WebsiteName: scan.Website.WebsiteName,

		Sha256: fmt.Sprintf("%x", utils.SHA256(scan.Website.WebsiteUrl)),
		SHA1:   fmt.Sprintf("%x", utils.SHA1(scan.Website.WebsiteUrl)),
		MD5:    fmt.Sprintf("%x", utils.MD5(scan.Website.WebsiteUrl)),

		Status: models.ScanStatusComplete,
	}

	// Create Scan
	scanResponse, err := controller.CreateScan(scanModel)
	if err != nil {
		fmt.Println("Failed to create scan", err)
		return models.Scan{}, err
	}

	// Create Certificates
	for _, certificate := range scan.Certificates {
		err := controller.CreateCertificate(scanResponse.ID, *certificate)
		if err != nil {
			fmt.Println("Failed to create certificate", err)
			return models.Scan{}, err
		}
	}

	// Create Findings
	controller.CreateFindings(scanResponse.ID, scan.Secrets)

	// Create Contents
	for _, script := range scan.Website.Scripts {
		content := models.Content{
			ScanID:  scanResponse.ID,
			Name:    script.Src,
			Content: script.Content,
		}

		_, err := controller.CreateContent(content)
		if err != nil {
			return models.Scan{}, err
		}
	}

	// Create Details
	detail := models.Detail{
		ScanID:       scanResponse.ID,
		IPAddresses:  scan.IPAddresses,
		HTTPHeaders:  scan.HTTPHeaders,
		IPRanges:     []string{},
		DNSNames:     []string{},
		PermittedDNS: []string{},
		ExcludedDNS:  []string{},
	}

	_, err = controller.CreateDetail(detail)
	if err != nil {
		return models.Scan{}, err
	}

	return scanResponse, nil
}
