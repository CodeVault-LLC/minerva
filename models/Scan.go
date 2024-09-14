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

	Sha256 string `gorm:"not null"`
	SHA1   string `gorm:"not null"`
	MD5    string `gorm:"not null"`

	Findings     []Finding     `gorm:"foreignKey:ScanID"`
	Contents     []Content     `gorm:"foreignKey:ScanID"`
	Certificates []Certificate `gorm:"foreignKey:ScanID"`
}

type ScanRequest struct {
	Url string `json:"url"`
}

type ScanResponse struct {
	WebsiteUrl  string `json:"websiteUrl"`
	WebsiteName string `json:"websiteName"`

	Scripts []ScriptRequest `json:"scripts"`
}

type ScanAPIResponse struct {
	ID uint `json:"id"`

	User UserMinimalResponse `json:"user"`

	WebsiteUrl  string `json:"website_url"`
	WebsiteName string `json:"website_name"`
	Status      string `json:"status"`

	Sha256 string `json:"sha256"`
	SHA1   string `json:"sha1"`
	MD5    string `json:"md5"`

	Findings     int64                 `json:"findings"`
	Certificates []CertificateResponse `json:"certificates"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
