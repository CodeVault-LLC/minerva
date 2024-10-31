package entities

import (
	"gorm.io/gorm"
)

type LicenseModel struct {
	gorm.Model

	License string `gorm:"not null"`

	Scans []ScanModel `gorm:"foreignKey:LicenseID"`
}
