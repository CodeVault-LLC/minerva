package models

import (
	"gorm.io/gorm"
)

type Scan struct {
	gorm.Model

	UserID uint
	User   User

	WebsiteUrl  string `gorm:"not null"`
	WebsiteName string `gorm:"not null"`
	Status      string `gorm:"not null"`

	Findings []Finding `gorm:"foreignKey:ScanID"`
	Contents []Content `gorm:"foreignKey:ScanID"`
}

type ScanRequest struct {
	WebsiteUrl  string `json:"websiteUrl"`
	WebsiteName string `json:"websiteName"`

	Scripts []ScriptRequest `json:"scripts"`
}

type ScanResponse struct {
	ID uint `json:"id"`

	User UserMinimalResponse `json:"user"`

	WebsiteUrl  string `json:"website_url"`
	WebsiteName string `json:"website_name"`

	Status string `json:"status"`

	Findings int64 `json:"findings"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
