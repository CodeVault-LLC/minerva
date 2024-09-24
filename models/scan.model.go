package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type ScanStatus string

const (
	ScanStatusArchived ScanStatus = "archived"
	ScanStatusComplete ScanStatus = "complete"
	ScanStatusFailed   ScanStatus = "failed"
	ScanStatusPending  ScanStatus = "pending"
)

type ScanModel struct {
	gorm.Model

	WebsiteUrl    string         `gorm:"not null"`
	WebsiteName   string         `gorm:"not null"`
	RedirectChain pq.StringArray `gorm:"type:text[]"`
	StatusCode    int            `gorm:"not null"`

	Status ScanStatus `gorm:"not null" default:"complete"`

	Sha256 string `gorm:"not null"`
	SHA1   string `gorm:"not null"`
	MD5    string `gorm:"not null"`

	LicenseID uint          `gorm:"not null"`
	License   *LicenseModel `gorm:"foreignKey:LicenseID"`

	Network NetworkModel `gorm:"foreignKey:ScanID"`

	Lists    []ListModel    `gorm:"foreignKey:ScanID"`
	Findings []FindingModel `gorm:"foreignKey:ScanID"`
	Contents []ContentModel `gorm:"foreignKey:ScanID"`
}

type ScanRequest struct {
	Url string `json:"url"`
}

type ScanResponse struct {
	WebsiteUrl  string `json:"websiteUrl"`
	WebsiteName string `json:"websiteName"`
	StatusCode  int    `json:"statusCode"`
}

type ScanAPIResponse struct {
	ID uint `json:"id"`

	WebsiteUrl    string         `json:"website_url"`
	WebsiteName   string         `json:"website_name"`
	RedirectChain pq.StringArray `json:"redirect_chain"`
	StatusCode    int            `json:"status_code"`

	Status string `json:"status"`

	Sha256 string `json:"sha256"`
	SHA1   string `json:"sha1"`
	MD5    string `json:"md5"`

	Findings int64          `json:"findings"`
	Lists    []ListResponse `json:"lists"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func ConvertScan(scan ScanModel) ScanAPIResponse {
	return ScanAPIResponse{
		ID:       scan.ID,
		Findings: int64(len(scan.Findings)),

		WebsiteUrl:    scan.WebsiteUrl,
		WebsiteName:   scan.WebsiteName,
		RedirectChain: scan.RedirectChain,
		StatusCode:    scan.StatusCode,

		Status: string(scan.Status),
		Sha256: scan.Sha256,
		SHA1:   scan.SHA1,
		MD5:    scan.MD5,

		Lists:     ConvertLists(scan.Lists),
		CreatedAt: scan.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: scan.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func ConvertScans(scans []ScanModel) []ScanAPIResponse {
	var scanResponses []ScanAPIResponse

	for _, scan := range scans {
		scanResponses = append(scanResponses, ConvertScan(scan))
	}

	return scanResponses
}
