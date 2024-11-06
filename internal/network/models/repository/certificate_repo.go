package repository

import (
	"crypto/x509"
	"encoding/base64"
	"encoding/json"

	"github.com/codevault-llc/minerva/internal/database"
	"github.com/codevault-llc/minerva/internal/network/models/entities"
	"github.com/codevault-llc/minerva/pkg/logger"
	"github.com/codevault-llc/minerva/pkg/utils"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"go.uber.org/zap"
)

type CertificateRepo struct {
	db *sqlx.DB
}

func NewCertificateRepository(db *sqlx.DB) *CertificateRepo {
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
		IsCa:                        cert.IsCA,
		DnsNames:                    pq.StringArray(cert.DNSNames),
		EmailAddresses:              pq.StringArray(cert.EmailAddresses),
		IpAddresses:                 pq.StringArray(utils.IPsToStrings(cert.IPAddresses)),
		Uris:                        pq.StringArray(utils.URIsToStrings(cert.URIs)),
		PermittedDnsDomainsCritical: cert.PermittedDNSDomainsCritical,
		PermittedDnsDomains:         pq.StringArray(cert.PermittedDNSDomains),
		ExcludedDnsDomains:          pq.StringArray(cert.ExcludedDNSDomains),
		PermittedIpRanges:           pq.StringArray(utils.IPNetsToStrings(cert.PermittedIPRanges)),
		ExcludedIpRanges:            pq.StringArray(utils.IPNetsToStrings(cert.ExcludedIPRanges)),
		PermittedEmailAddresses:     pq.StringArray(cert.PermittedEmailAddresses),
		ExcludedEmailAddresses:      pq.StringArray(cert.ExcludedEmailAddresses),
		PermittedUriDomains:         pq.StringArray(cert.PermittedURIDomains),
		ExcludedUriDomains:          pq.StringArray(cert.ExcludedURIDomains),
	}

	query, values, err := database.StructToQuery(certificate, "certificates")
	if err != nil {
		logger.Log.Error("Failed to generate query", zap.Error(err))
		return entities.CertificateModel{}, err
	}

	tx, err := repository.db.Beginx()
	if err != nil {
		logger.Log.Error("Failed to start transaction", zap.Error(err))
		return entities.CertificateModel{}, err
	}

	_, err = database.InsertStruct(tx, query, values)
	if err != nil {
		logger.Log.Error("Failed to insert certificate", zap.Error(err))
		err := tx.Rollback()

		if err != nil {
			logger.Log.Error("Failed to rollback transaction", zap.Error(err))
		}
		return entities.CertificateModel{}, err
	}

	err = tx.Commit()
	if err != nil {
		logger.Log.Error("Failed to commit transaction", zap.Error(err))
		return entities.CertificateModel{}, err
	}

	return certificate, nil
}
