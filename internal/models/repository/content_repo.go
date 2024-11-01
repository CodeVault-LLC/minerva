package repository

import (
	"github.com/codevault-llc/humblebrag-api/internal/models/entities"
	"gorm.io/gorm"
)

type ContentRepo struct {
	db *gorm.DB
}

func NewContentRepo(db *gorm.DB) *ContentRepo {
	return &ContentRepo{db: db}
}

var ContentRepository *ContentRepo

func (repository *ContentRepo) SaveContentResult(content entities.ContentModel) (entities.ContentModel, error) {
	tx := repository.db.Begin()
	if err := tx.Create(&content).Error; err != nil {
		tx.Rollback()
		return entities.ContentModel{}, err
	}

	tx.Commit()
	return content, nil
}

func (repository *ContentRepo) FindContentByHash(hashedBody string) (entities.ContentModel, error) {
	var content entities.ContentModel
	if err := repository.db.Where("hashed_body = ?", hashedBody).First(&content).Error; err != nil {
		return content, err
	}

	return content, nil
}

func (repository *ContentRepo) IncrementAccessCount(contentID uint) error {
	tx := repository.db.Begin()
	if err := tx.Model(&entities.ContentModel{}).Where("id = ?", contentID).Update("access_count", gorm.Expr("access_count + ?", 1)).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (repository *ContentRepo) AddContentToScan(scanID uint, contentID uint) error {
	tx := repository.db.Begin()
	var scan entities.ScanModel
	if err := tx.First(&scan, scanID).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&scan).Association("Contents").Append(&entities.ContentModel{ID: contentID}); err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (repository *ContentRepo) CreateContentStorage(storage entities.ContentStorageModel) error {
	tx := repository.db.Begin()
	if err := tx.Create(&storage).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
