package models

import (
	"time"

	"gorm.io/gorm"
)

type Certificate struct {
	gorm.Model

	ScanID uint
	Scan   Scan

	Subject string `gorm:"not null"`
	Issuer  string `gorm:"not null"`

	NotBefore time.Time `gorm:"not null"`
	NotAfter  time.Time `gorm:"not null"`

	SignatureAlgorithm string `gorm:"not null"`
	PublicKeyAlgorithm string `gorm:"not null"`
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
