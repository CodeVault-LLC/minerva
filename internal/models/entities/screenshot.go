package entities

import "gorm.io/gorm"

type ScreenshotModel struct {
	gorm.Model

	RedirectId uint
	Redirect   *RedirectModel

	ImageBucket    string `gorm:"type:text"`
	ImageObjectKey string `gorm:"type:text"`

	CompressedSize int `gorm:"type:integer"`
}
