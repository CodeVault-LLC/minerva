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

	// HTTP fields
	HTTPHeaders pq.StringArray `gorm:"type:text[]"` // PostgreSQL array

	// Relationships
	Whois        WhoisModel         `gorm:"foreignKey:NetworkId"`
	Certificates []CertificateModel `gorm:"foreignKey:NetworkId"` // Separate foreign key for certificates
	DNS          DNSModel           `gorm:"foreignKey:NetworkId"`
}

type NetworkResponse struct {
	ID uint `json:"id"`

	IPAddresses []string `json:"ip_addresses"`
	IPRanges    []string `json:"ip_ranges"`

	HTTPHeaders []string `json:"http_headers"`

	Certificates []CertificateResponse `json:"certificates"`
	Whois        WhoisResponse         `json:"whois"`
	DNS          DNSResponse           `json:"dns"`
}

func ConvertNetwork(network NetworkModel) NetworkResponse {
	convertedCertificates := make([]CertificateResponse, len(network.Certificates))

	for i, cert := range network.Certificates {
		convertedCertificates[i] = ConvertCertificate(cert)
	}

	return NetworkResponse{
		ID:          network.ID,
		IPAddresses: network.IPAddresses,
		IPRanges:    network.IPRanges,
		HTTPHeaders: network.HTTPHeaders,

		Certificates: convertedCertificates,
		Whois:        ConvertWhois(network.Whois),
		DNS:          ConvertDNS(network.DNS),
	}
}
