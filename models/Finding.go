package models

import "gorm.io/gorm"

type Finding struct {
	gorm.Model

	ScanID uint
	Scan   Scan

	RegexName        string `gorm:"not null"`
	RegexDescription string `gorm:"not null"`

	Match  string `gorm:"not null"`
	Source string `gorm:"not null"`
	Line   int    `gorm:"not null"`
}

type ScriptRequest struct {
	Src     string `json:"src"`
	Content string `json:"content"`
}

type FindingResponse struct {
	ID     uint `json:"id"`
	ScanID uint `json:"scan_id"`

	RegexName        string `json:"regex_name"`
	RegexDescription string `json:"regex_description"`

	Match  string `json:"match"`
	Source string `json:"source"`
	Line   int    `json:"line"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
