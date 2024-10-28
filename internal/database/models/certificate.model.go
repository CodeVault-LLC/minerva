package models

import (
	"crypto/x509"
	"encoding/json"
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type CertificateModel struct {
	gorm.Model

	NetworkId uint
	Network   NetworkModel `gorm:"foreignKey:NetworkId"`

	Subject string `gorm:"not null"`
	Issuer  string `gorm:"not null"`

	NotBefore time.Time `gorm:"not null"`
	NotAfter  time.Time `gorm:"not null"`

	SignatureAlgorithm x509.SignatureAlgorithm `gorm:"not null"`
	Signature          []byte                  `gorm:"not null"`

	PublicKeyAlgorithm x509.PublicKeyAlgorithm `gorm:"not null"`
	PublicKey          string                  `gorm:"not null"` // Store as JSON

	SerialNumber string        `gorm:"not null"` // Store as string
	Version      int           `gorm:"not null"`
	KeyUsage     x509.KeyUsage `gorm:"not null"`

	BasicConstraintsValid bool `gorm:"not null"`
	IsCA                  bool `gorm:"not null"`

	DNSNames       pq.StringArray `gorm:"type:text[]"` // PostgreSQL array
	EmailAddresses pq.StringArray `gorm:"type:text[]"` // PostgreSQL array
	IPAddresses    pq.StringArray `gorm:"type:text[]"` // Store as string representation of IPs
	URIs           pq.StringArray `gorm:"type:text[]"` // Store as string representation of URLs

	PermittedDNSDomainsCritical bool           `gorm:"not null"`
	PermittedDNSDomains         pq.StringArray `gorm:"type:text[]"` // PostgreSQL array
	ExcludedDNSDomains          pq.StringArray `gorm:"type:text[]"` // PostgreSQL array
	PermittedIPRanges           pq.StringArray `gorm:"type:text[]"` // Store as string representation of IP ranges
	ExcludedIPRanges            pq.StringArray `gorm:"type:text[]"` // Store as string representation of IP ranges
	PermittedEmailAddresses     pq.StringArray `gorm:"type:text[]"` // PostgreSQL array
	ExcludedEmailAddresses      pq.StringArray `gorm:"type:text[]"` // PostgreSQL array
	PermittedURIDomains         pq.StringArray `gorm:"type:text[]"` // PostgreSQL array
	ExcludedURIDomains          pq.StringArray `gorm:"type:text[]"` // PostgreSQL array
}

// Custom methods to handle PublicKey marshaling/unmarshaling
func (c *CertificateModel) BeforeSave(tx *gorm.DB) (err error) {
	if c.PublicKey != "" {
		encodedPublicKey, err := json.Marshal(c.PublicKey)
		if err != nil {
			return err
		}
		c.PublicKey = string(encodedPublicKey)
	}
	return nil
}

func (c *CertificateModel) AfterFind(tx *gorm.DB) (err error) {
	if len(c.PublicKey) > 0 {
		err = json.Unmarshal([]byte(c.PublicKey), &c.PublicKey)
		if err != nil {
			return err
		}
	}
	return nil
}

type CertificateResponse struct {
	ID uint `json:"id"`

	Subject string `json:"subject"`
	Issuer  string `json:"issuer"`

	NotBefore time.Time `json:"not_before"`
	NotAfter  time.Time `json:"not_after"`

	SignatureAlgorithm string `json:"signature_algorithm"`
	PublicKeyAlgorithm string `json:"public_key_algorithm"`
}

func ConvertCertificate(certificate CertificateModel) CertificateResponse {
	return CertificateResponse{
		ID:                 certificate.ID,
		Issuer:             certificate.Issuer,
		Subject:            certificate.Subject,
		NotBefore:          certificate.NotBefore,
		NotAfter:           certificate.NotAfter,
		SignatureAlgorithm: certificate.SignatureAlgorithm.String(),
		PublicKeyAlgorithm: certificate.PublicKeyAlgorithm.String(),
	}
}

func ConvertCertificates(certificates []CertificateModel) []CertificateResponse {
	var certificateResponses []CertificateResponse

	for _, certificate := range certificates {
		certificateResponses = append(certificateResponses, CertificateResponse{
			ID:                 certificate.ID,
			Issuer:             certificate.Issuer,
			Subject:            certificate.Subject,
			NotBefore:          certificate.NotBefore,
			NotAfter:           certificate.NotAfter,
			SignatureAlgorithm: certificate.SignatureAlgorithm.String(),
			PublicKeyAlgorithm: certificate.PublicKeyAlgorithm.String(),
		})
	}

	if len(certificateResponses) == 0 {
		return []CertificateResponse{}
	}

	return certificateResponses
}
