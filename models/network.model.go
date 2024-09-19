package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type NetworkModel struct {
	gorm.Model

	ScanID uint
	Scan   *ScanModel

	// Detail fields
	IPAddresses pq.StringArray `gorm:"type:text[]"` // PostgreSQL array
	IPRanges    pq.StringArray `gorm:"type:text[]"` // PostgreSQL array

	// DNS fields
	DNSNames     pq.StringArray `gorm:"type:text[]"` // PostgreSQL array
	PermittedDNS pq.StringArray `gorm:"type:text[]"` // PostgreSQL array
	ExcludedDNS  pq.StringArray `gorm:"type:text[]"` // PostgreSQL array

	// HTTP fields
	HTTPHeaders pq.StringArray `gorm:"type:text[]"` // PostgreSQL array

	// Relationships
	Whois        WhoisModel         `gorm:"foreignKey:ScanID"`
	Certificates []CertificateModel `gorm:"foreignKey:ScanID"`
}

type NetworkResponse struct {
	ID uint `json:"id"`

	IPAddresses []string `json:"ip_addresses"`
	IPRanges    []string `json:"ip_ranges"`

	DNSNames     []string `json:"dns_names"`
	PermittedDNS []string `json:"permitted_dns"`
	ExcludedDNS  []string `json:"excluded_dns"`

	HTTPHeaders []string `json:"http_headers"`
}

func ConvertNetwork(network NetworkModel) NetworkResponse {
	return NetworkResponse{
		ID:           network.ID,
		IPAddresses:  network.IPAddresses,
		IPRanges:     network.IPRanges,
		DNSNames:     network.DNSNames,
		PermittedDNS: network.PermittedDNS,
		ExcludedDNS:  network.ExcludedDNS,
		HTTPHeaders:  network.HTTPHeaders,
	}
}
