package repository

import (
	"github.com/codevault-llc/humblebrag-api/internal/database"
	"github.com/codevault-llc/humblebrag-api/internal/models/entities"
	"github.com/codevault-llc/humblebrag-api/internal/models/viewmodels"
	"github.com/codevault-llc/humblebrag-api/pkg/logger"
	"go.uber.org/zap"
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

func (repository *ContentRepo) GetScanContents(scanID uint) ([]viewmodels.Contents, error) {
	var scan entities.ScanModel

	// Retrieve the scan by ID, preloading the associated contents.
	if err := database.DB.Preload("Contents").First(&scan, scanID).Error; err != nil {
		return nil, err
	}

	// Extract the associated contents.
	content := scan.Contents

	logger.Log.Info("Retrieved contents for scan", zap.Uint("scanID", scanID), zap.Int("contentCount", len(content)))

	// Create maps to hold associated tags and storage information.
	contentIDs := make([]uint, len(content))
	for i, c := range content {
		contentIDs[i] = c.ID
	}

	// Retrieve associated tags for each content ID.
	tagsMap := make(map[uint][]string)
	var tags []entities.ContentTagsModel
	if err := database.DB.Where("content_id IN ?", contentIDs).Find(&tags).Error; err != nil {
		return nil, err
	}

	for _, tag := range tags {
		tagsMap[tag.ContentID] = append(tagsMap[tag.ContentID], tag.Tag)
	}

	// Retrieve associated storage information for each content ID.
	storageMap := make(map[uint]entities.ContentStorageModel)
	var storageRecords []entities.ContentStorageModel
	if err := database.DB.Where("content_id IN ?", contentIDs).Find(&storageRecords).Error; err != nil {
		return nil, err
	}

	for _, storageRecord := range storageRecords {
		storageMap[storageRecord.ContentID] = storageRecord
	}

	// Convert the content models into the content responses with tags and storage details.
	return viewmodels.ConvertContents(content, tagsMap, storageMap), nil
}

func (repository *ContentRepo) GetScanContent(contentID uint) (entities.ContentModel, error) {
	var content entities.ContentModel

	// Retrieve the content by ID, preloading the associated tags and storage information.
	if err := database.DB.Preload("Tags").Preload("Storage").First(&content, contentID).Error; err != nil {
		return content, err
	}

	return content, nil
}
