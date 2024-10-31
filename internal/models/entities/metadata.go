package entities

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type MetadataModel struct {
	gorm.Model

	ScanID uint
	Scan   *ScanModel `gorm:"foreignKey:ScanID"`

	// Detail fields
	Robots  string `gorm:"not null"`
	Readme  string `gorm:"not null"`
	License string `gorm:"not null"`

	// More special detailed into CMS
	CMS            string         `gorm:"not null"`
	ServerSoftware string         `gorm:"not null"`
	Frameworks     pq.StringArray `gorm:"type:text[]"` // PostgreSQL array
	ServerLanguage string         `gorm:"not null"`
}
