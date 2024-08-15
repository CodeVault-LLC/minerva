package models

import "gorm.io/gorm"

type Content struct {
	gorm.Model

	ScanID uint
	Scan   Scan

	Name string `gorm:"not null"`

	Content string `gorm:"not null"`
}

type ContentResponse struct {
	ID      uint   `json:"id"`
	ScanID  uint   `json:"scan_id"`
	Name    string `json:"name"`
	Content string `json:"content"`
}
