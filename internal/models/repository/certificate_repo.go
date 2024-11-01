package repository

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/codevault-llc/humblebrag-api/internal/models/entities"
	"github.com/codevault-llc/humblebrag-api/pkg/logger"
	"github.com/codevault-llc/humblebrag-api/pkg/utils"
	"github.com/lib/pq"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type CertificateRepo struct {
	db *gorm.DB
}

func NewCertificateRepository(db *gorm.DB) *CertificateRepo {
	return &CertificateRepo{
		db: db,
	}
}

var CertificateRepository *CertificateRepo

func (repository *CertificateRepo) Create(networkId uint, cert x509.Certificate) (entities.CertificateModel, error) {
	publicKeyJSON, err := json.Marshal(cert.PublicKey)
	if err != nil {
		return entities.CertificateModel{}, err
	}

	signatureBase64 := base64.StdEncoding.EncodeToString(cert.Signature)

	certificate := entities.CertificateModel{
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

	fmt.Println("Certificate continue its path")

	tx := repository.db.Begin()
	if err := tx.Create(&certificate).Error; err != nil {
		tx.Rollback()
		logger.Log.Error("Failed to create certificate", zap.Error(err))
		return entities.CertificateModel{}, err
	}

	logger.Log.Info("Created certificate")
	tx.Commit()
	return certificate, nil
}
