package service

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/json"

	"github.com/codevault-llc/humblebrag-api/constants"
	"github.com/codevault-llc/humblebrag-api/models"
	"github.com/codevault-llc/humblebrag-api/utils"
	"github.com/lib/pq"
)

func CreateScan(scan models.Scan) (models.Scan, error) {
	constants.DB.Model(&models.Scan{}).Where("website_url = ?", scan.WebsiteUrl).Update("status", models.ScanStatusArchived)

	if err := constants.DB.Create(&scan).Error; err != nil {
		return scan, err
	}

	return scan, nil
}

func CreateFindings(scanID uint, secrets []utils.RegexReturn) {
	for _, secret := range secrets {
		for _, match := range secret.Matches {
			finding := models.Finding{
				ScanID: scanID,
				Line:   match.Line,
				Match:  match.Match,
				Source: match.Source,

				RegexName:        secret.Name,
				RegexDescription: secret.Description,
			}

			constants.DB.Create(&finding)
		}
	}
}

func CreateCertificate(scanID uint, cert x509.Certificate) error {
	publicKeyJSON, err := json.Marshal(cert.PublicKey)
	if err != nil {
		return err
	}

	signatureBase64 := base64.StdEncoding.EncodeToString(cert.Signature)

	certificate := models.Certificate{
		ScanID:                      scanID,
		Issuer:                      string(cert.Issuer.CommonName),
		Subject:                     string(cert.Subject.CommonName),
		NotBefore:                   cert.NotBefore,
		NotAfter:                    cert.NotAfter,
		SignatureAlgorithm:          cert.SignatureAlgorithm,
		Signature:                   []byte(signatureBase64),
		PublicKeyAlgorithm:          cert.PublicKeyAlgorithm,
		PublicKey:                   string(publicKeyJSON),
		SerialNumber:                cert.SerialNumber.String(),
		Version:                     cert.Version,
		KeyUsage:                    cert.KeyUsage,
		BasicConstraintsValid:       cert.BasicConstraintsValid,
		IsCA:                        cert.IsCA,
		DNSNames:                    pq.StringArray(cert.DNSNames),
		EmailAddresses:              pq.StringArray(cert.EmailAddresses),
		IPAddresses:                 pq.StringArray(utils.IPsToStrings(cert.IPAddresses)),
		URIs:                        pq.StringArray(utils.URIsToStrings(cert.URIs)),
		PermittedDNSDomainsCritical: cert.PermittedDNSDomainsCritical,
		PermittedDNSDomains:         pq.StringArray(cert.PermittedDNSDomains),
		ExcludedDNSDomains:          pq.StringArray(cert.ExcludedDNSDomains),
		PermittedIPRanges:           pq.StringArray(utils.IPNetsToStrings(cert.PermittedIPRanges)),
		ExcludedIPRanges:            pq.StringArray(utils.IPNetsToStrings(cert.ExcludedIPRanges)),
		PermittedEmailAddresses:     pq.StringArray(cert.PermittedEmailAddresses),
		ExcludedEmailAddresses:      pq.StringArray(cert.ExcludedEmailAddresses),
		PermittedURIDomains:         pq.StringArray(cert.PermittedURIDomains),
		ExcludedURIDomains:          pq.StringArray(cert.ExcludedURIDomains),
	}

	constants.DB.Create(&certificate)
	return nil
}

func CreateContent(content models.Content) (models.Content, error) {
	if err := constants.DB.Create(&content).Error; err != nil {
		return content, err
	}

	return content, nil
}

func CreateDetail(detail models.Detail) (models.Detail, error) {
	if err := constants.DB.Create(&detail).Error; err != nil {
		return detail, err
	}

	return detail, nil
}

func CreateList(list models.List) (models.List, error) {
	if err := constants.DB.Create(&list).Error; err != nil {
		return list, err
	}

	return list, nil
}

func CreateLists(lists []models.List) ([]models.List, error) {
	if err := constants.DB.Create(&lists).Error; err != nil {
		return lists, err
	}

	return lists, nil
}

func GetScans() ([]models.ScanAPIResponse, error) {
	var scans []models.Scan

	if err := constants.DB.Preload("Findings").Where("status = ?", models.ScanStatusComplete).Find(&scans).Error; err != nil {
		return models.ConvertScans(scans), err
	}

	return models.ConvertScans(scans), nil
}

func GetScan(scanID string) (models.ScanAPIResponse, error) {
	var scan models.Scan

	if err := constants.DB.Where("id = ?", scanID).Preload("Findings").Preload("Lists").Preload("Detail").Preload("Certificates").Where("status = ?", models.ScanStatusComplete).First(&scan).Error; err != nil {
		return models.ConvertScan(scan), err
	}

	return models.ConvertScan(scan), nil
}

func GetScanFindings(scanID string) ([]models.FindingResponse, error) {
	var findings []models.Finding

	if err := constants.DB.Where("scan_id = ?", scanID).
		Find(&findings).
		Error; err != nil {
		return models.ConvertFindings(findings), err
	}

	return models.ConvertFindings(findings), nil
}

func GetScanContent(scanID string) ([]models.ContentResponse, error) {
	var content []models.Content

	if err := constants.DB.Where("scan_id = ?", scanID).
		Find(&content).
		Error; err != nil {
		return models.ConvertContents(content), err
	}

	return models.ConvertContents(content), nil
}