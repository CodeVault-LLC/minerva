package models

import "gorm.io/gorm"

type List struct {
	gorm.Model

	ScanID uint
	Scan   Scan

	ListID string // match towards the listID in the config
}

type ListResponse struct {
	ID uint `json:"id"`

	Description string   `json:"description"`
	ListID      string   `json:"list_id"`
	Categories  []string `json:"categories"`
	URL         string   `json:"url"`
}
