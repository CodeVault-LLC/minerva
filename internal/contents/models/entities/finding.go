package entities

import (
	"github.com/codevault-llc/humblebrag-api/internal/models/entities"
	"gorm.io/gorm"
)

type FindingModel struct {
	gorm.Model

	ScanID uint
	Scan   entities.ScanModel `gorm:"foreignKey:ScanID"`

	RegexName        string `gorm:"not null"`
	RegexDescription string `gorm:"not null"`

	Match  string `gorm:"not null"`
	Source string `gorm:"not null"`
	Line   int    `gorm:"not null"`
}
