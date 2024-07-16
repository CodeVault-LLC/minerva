package models

import "gorm.io/gorm"

type Script struct {
	gorm.Model

	ScanID uint
	Scan Scan

	Src string `gorm:"not null"`
	Status string `gorm:"not null"`
	Content string `gorm:"not null"`

	Secrets []Secret `gorm:"foreignKey:ScriptID"`
}

type ScriptRequest struct {
	Src string `json:"src"`
	Content string `json:"content"`
}

type ScriptResponse struct {
	ID uint `json:"id"`
	ScanID uint `json:"scan_id"`
	Src string `json:"src"`
	Status string `json:"status"`
	Content string `json:"content"`
}
