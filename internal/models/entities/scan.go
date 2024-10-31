package entities

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type ScanModel struct {
	gorm.Model

	Url        string `gorm:"not null"`
	Title      string `gorm:"not null"`
	StatusCode int    `gorm:"not null"`

	RedirectChain pq.StringArray `gorm:"type:text[]"`

	Status ScanStatus `gorm:"not null;default:'complete'"`

	Sha256 string `gorm:"not null"`
	SHA1   string `gorm:"not null"`
	MD5    string `gorm:"not null"`

	LicenseID uint          `gorm:"not null"`
	License   *LicenseModel `gorm:"foreignKey:LicenseID"`

	Network  NetworkModel  `gorm:"foreignKey:ScanID"`
	Metadata MetadataModel `gorm:"foreignKey:ScanID"`

	Lists    []FilterModel  `gorm:"foreignKey:ScanID"`
	Findings []FindingModel `gorm:"foreignKey:ScanID"`

	// Define the many-to-many relationship through the join table.
	Contents []ContentModel `gorm:"many2many:scan_contents"`
}

type ScanStatus string

const (
	ScanStatusArchived ScanStatus = "archived"
	ScanStatusComplete ScanStatus = "complete"
	ScanStatusFailed   ScanStatus = "failed"
	ScanStatusPending  ScanStatus = "pending"
)
