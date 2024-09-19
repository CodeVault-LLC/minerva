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

func CreateCertificate(networkId uint, cert x509.Certificate) error {
	publicKeyJSON, err := json.Marshal(cert.PublicKey)
	if err != nil {
		return err
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

	constants.DB.Create(&certificate)
	return nil
}
