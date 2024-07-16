package models

import (
	"gorm.io/gorm"
)

type Scan struct {
	gorm.Model

	UserID uint
	User User

	Url string `gorm:"not null"`
	Status string `gorm:"not null"`

	Scripts []Script `gorm:"foreignKey:ScanID"`
}

type ScanRequest struct {
	Url string `json:"url"`
	Depth int `json:"depth"`

	DoScripts bool `json:"doScripts"`
	DoLinks bool `json:"doLinks"`
	DoImages bool `json:"doImages"`
	DoStyles bool `json:"doStyles"`

	Scripts []ScriptRequest `json:"scripts"`
}

type ScanResponse struct {
	ID uint `json:"id"`
	UserID uint `json:"user_id"`
	Url string `json:"url"`
	Status string `json:"status"`
}
