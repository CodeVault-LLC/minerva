package entities

import (
	"time"

	"gorm.io/gorm"
)

type RedirectModel struct {
	gorm.Model

	ScanID uint
	Scan   *ScanModel

	Url        string `gorm:"type:text"`
	HttpStatus int    `gorm:"type:integer"`

	Timestamp time.Time `gorm:"type:timestamp"`

	// Relationships
	Screenshot ScreenshotModel `gorm:"foreignKey:RedirectId"`
}
