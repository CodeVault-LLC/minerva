package models

import (
	"gorm.io/gorm"
)

type Detail struct {
	gorm.Model

	// ScanID is the foreign key for the Scan model
	ScanID uint
	Scan   Scan

	// Detail fields
	IPAddresses []string `gorm:"type:text[]"` // Store as string representation of IPs
	IPRanges    []string `gorm:"type:text[]"` // Store as string representation of IP ranges

	// DNS fields
	DNSNames     []string `gorm:"type:text[]"` // PostgreSQL array
	PermittedDNS []string `gorm:"type:text[]"` // PostgreSQL array
	ExcludedDNS  []string `gorm:"type:text[]"` // PostgreSQL array

	// HTTP fields
	HTTPHeaders map[string][]string `gorm:"type:jsonb"` // Store as JSONB
}
