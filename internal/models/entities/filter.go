package entities

import "gorm.io/gorm"

type FilterModel struct {
	gorm.Model

	ScanID uint
	Scan   ScanModel

	FilterID string `gorm:"not null"`
}
