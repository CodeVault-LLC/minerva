package service

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/json"

	"github.com/codevault-llc/humblebrag-api/internal/database"
	"github.com/codevault-llc/humblebrag-api/models"
	"github.com/codevault-llc/humblebrag-api/pkg/utils"
	"github.com/lib/pq"
)

func CreateCertificate(networkId uint, cert x509.Certificate) (models.CertificateModel, error) {
	publicKeyJSON, err := json.Marshal(cert.PublicKey)
	if err != nil {
		return models.CertificateModel{}, err
	}

	signatureBase64 := base64.StdEncoding.EncodeToString(cert.Signature)

	certificate := models.CertificateModel{
		NetworkId:                   networkId,
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

	database.DB.Create(&certificate)
	return certificate, nil
}

func DeleteCertificates(networkId uint) error {
	if err := database.DB.Where("network_id = ?", networkId).Delete(&models.CertificateModel{}).Error; err != nil {
		return err
	}

	return nil
}

func CreateCertificateResult(certId uint, result models.CertificateResult) error {
	tx := database.DB.Begin()

	certResult := models.CertificateResultModel{
		CertificateId:   certId,
		Expired:         result.Expired,
		Trusted:         result.Trusted,
		Weak:            result.Weak,
		MaliciousIssuer: result.MaliciousIssuer,
		Revoked:         result.Revoked,
		Domain:          result.Domain,
	}

	if err := tx.Create(&certResult).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
