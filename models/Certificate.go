package models

import "gorm.io/gorm"

type Certificate struct {
	gorm.Model

	ScanID uint
	Scan   Scan

	Subject string `gorm:"not null"`
	Issuer  string `gorm:"not null"`
	Valid   bool   `gorm:"not null"`
}

