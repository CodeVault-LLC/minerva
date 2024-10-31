package entities

import "gorm.io/gorm"

type FindingModel struct {
	gorm.Model

	ScanID uint
	Scan   ScanModel

	RegexName        string `gorm:"not null"`
	RegexDescription string `gorm:"not null"`

	Match  string `gorm:"not null"`
	Source string `gorm:"not null"`
	Line   int    `gorm:"not null"`
}
