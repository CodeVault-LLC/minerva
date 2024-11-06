package viewmodels

import "github.com/codevault-llc/minerva/internal/core/models/entities"

type Scan struct {
	ID uint `json:"id"`

	Url        string `json:"url"`
	Title      string `json:"title"`
	StatusCode int    `json:"status_code"`

	Status string `json:"status"`

	Sha256 string `json:"sha256"`
	SHA1   string `json:"sha1"`
	MD5    string `json:"md5"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func ConvertScan(scan entities.ScanModel) Scan {
	return Scan{
		ID: uint(scan.Id),

		Url:        scan.Url,
		Title:      scan.Title,
		StatusCode: scan.StatusCode,

		Status: string(scan.Status),
		Sha256: scan.Sha256,
		SHA1:   scan.Sha1,
		MD5:    scan.Md5,

		CreatedAt: scan.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: scan.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func ConvertScans(scans []entities.ScanModel) []Scan {
	var scanResponses []Scan

	for _, scan := range scans {
		scanResponses = append(scanResponses, ConvertScan(scan))
	}

	return scanResponses
}
