package repository

import (
	"github.com/codevault-llc/humblebrag-api/internal/models/entities"
	"gorm.io/gorm"
)

type ScreenshotRepo struct {
	db *gorm.DB
}

func NewScreenshotRepo(db *gorm.DB) *ScreenshotRepo {
	return &ScreenshotRepo{db: db}
}

var ScreenshotRepository *ScreenshotRepo

func (r *ScreenshotRepo) Create(screenshot entities.ScreenshotModel) (entities.ScreenshotModel, error) {
	tx := r.db.Begin()
	if err := tx.Create(&screenshot).Error; err != nil {
		tx.Rollback()
		return entities.ScreenshotModel{}, err
	}
	tx.Commit()
	return screenshot, nil
}
