package models

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
	ServerLanguage       string         `gorm:"not null"`
}

type MetadataResponse struct {
	ID uint `json:"id"`

	Robots  string `json:"robots"`
	Readme  string `json:"readme"`
	License string `json:"license"`

	CMS            string   `json:"cms"`
	ServerSoftware string   `json:"server_software"`
	Frameworks     []string `json:"frameworks"`
	ServerLanguage       string   `json:"server_language"`
}

func ConvertMetadata(metadata MetadataModel) MetadataResponse {
	return MetadataResponse{
		ID: metadata.ID,

		Robots:  metadata.Robots,
		Readme:  metadata.Readme,
		License: metadata.License,

		CMS:            metadata.CMS,
		ServerSoftware: metadata.ServerSoftware,
		Frameworks:     metadata.Frameworks,
	}
}
