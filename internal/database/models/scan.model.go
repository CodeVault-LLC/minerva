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

type ScanRequest struct {
	Url       string `json:"url"`
	UserAgent string `json:"userAgent"`
}

type ScanAPIResponse struct {
	ID uint `json:"id"`

	Url        string `json:"url"`
	Title      string `json:"title"`
	StatusCode int    `json:"status_code"`

	RedirectChain pq.StringArray `json:"redirect_chain"`

	Status string `json:"status"`

	Sha256 string `json:"sha256"`
	SHA1   string `json:"sha1"`
	MD5    string `json:"md5"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func ConvertScan(scan ScanModel) ScanAPIResponse {
	return ScanAPIResponse{
		ID: scan.ID,

		Url:           scan.Url,
		Title:         scan.Title,
		RedirectChain: scan.RedirectChain,
		StatusCode:    scan.StatusCode,

		Status: string(scan.Status),
		Sha256: scan.Sha256,
		SHA1:   scan.SHA1,
		MD5:    scan.MD5,

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
