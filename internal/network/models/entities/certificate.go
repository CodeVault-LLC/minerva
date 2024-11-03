package entities

import (
	"crypto/x509"
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
