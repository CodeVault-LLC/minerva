package models

import (
	"gorm.io/gorm"
)

type ScanStatus string

const (
	ScanStatusArchived ScanStatus = "archived"
	ScanStatusComplete ScanStatus = "complete"
	ScanStatusFailed   ScanStatus = "failed"
)

type Scan struct {
	gorm.Model

	WebsiteUrl  string     `gorm:"not null"`
	WebsiteName string     `gorm:"not null"`
	Status      ScanStatus `gorm:"not null" default:"complete"`

	Sha256 string `gorm:"not null"`
	SHA1   string `gorm:"not null"`
	MD5    string `gorm:"not null"`

	UserID uint
	User   User

	Detail       Detail        `gorm:"foreignKey:ScanID"`
	Lists        []List        `gorm:"foreignKey:ScanID"`
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

	WebsiteUrl  string `json:"website_url"`
	WebsiteName string `json:"website_name"`
	Status      string `json:"status"`

	Sha256 string `json:"sha256"`
	SHA1   string `json:"sha1"`
	MD5    string `json:"md5"`

	Detail       DetailResponse        `json:"detail"`
	Findings     int64                 `json:"findings"`
	Certificates []CertificateResponse `json:"certificates"`
	Lists        []ListResponse        `json:"lists"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func ConvertScan(scan Scan) ScanAPIResponse {
	return ScanAPIResponse{
		ID:       scan.ID,
		Findings: int64(len(scan.Findings)),

		WebsiteUrl:  scan.WebsiteUrl,
		WebsiteName: scan.WebsiteName,
		Status:      string(scan.Status),
		Sha256:      scan.Sha256,
		SHA1:        scan.SHA1,
		MD5:         scan.MD5,

		Certificates: ConvertCertificates(scan.Certificates),
		Detail:       ConvertDetail(scan.Detail),
		Lists:        ConvertLists(scan.Lists),
		CreatedAt:    scan.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    scan.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func ConvertScans(scans []Scan) []ScanAPIResponse {
	var scanResponses []ScanAPIResponse

	for _, scan := range scans {
		scanResponses = append(scanResponses, ConvertScan(scan))
	}

	return scanResponses
}
