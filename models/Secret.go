package models

import "gorm.io/gorm"

type Secret struct {
	gorm.Model

	ScriptID uint
	Script Script

	Secret string `gorm:"not null"`
	Source string `gorm:"not null"`
	Line int `gorm:"not null"`
}

type SecretResponse struct {
	ID uint `json:"id"`
	ScriptID uint `json:"script_id"`
	Secret string
	Source string
	Line int
}
