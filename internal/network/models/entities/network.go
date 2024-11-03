package entities

import (
	"github.com/codevault-llc/humblebrag-api/internal/models/entities"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type NetworkModel struct {
	gorm.Model

	ScanID uint
	Scan   *entities.ScanModel `gorm:"foreignKey:ScanID"`

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
