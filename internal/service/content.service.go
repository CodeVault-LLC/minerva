package service

import (
	"github.com/codevault-llc/humblebrag-api/internal/database"
	"github.com/codevault-llc/humblebrag-api/internal/database/models"
	"github.com/codevault-llc/humblebrag-api/pkg/logger"
	"go.uber.org/zap"
)

// CreateContent creates content in the database
func CreateContent(content models.ContentModel) (models.ContentModel, error) {
	if err := database.DB.Create(&content).Error; err != nil {
		return content, err
	}

	return content, nil
}

func CreateContentTags(tags []models.ContentTagsModel) error {
	if err := database.DB.Create(&tags).Error; err != nil {
		return err
	}

	return nil
}

func CreateContentStorage(storage models.ContentStorageModel) error {
	if err := database.DB.Create(&storage).Error; err != nil {
		return err
	}

	return nil
}

// CreateContents gets content from the database
func GetScanContent(scanID uint) ([]models.ContentResponse, error) {
	var scan models.ScanModel

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
	var tags []models.ContentTagsModel
	if err := database.DB.Where("content_id IN ?", contentIDs).Find(&tags).Error; err != nil {
		return nil, err
	}

	for _, tag := range tags {
		tagsMap[tag.ContentID] = append(tagsMap[tag.ContentID], tag.Tag)
	}

	// Retrieve associated storage information for each content ID.
	storageMap := make(map[uint]models.ContentStorageModel)
	var storageRecords []models.ContentStorageModel
	if err := database.DB.Where("content_id IN ?", contentIDs).Find(&storageRecords).Error; err != nil {
		return nil, err
	}

	for _, storageRecord := range storageRecords {
		storageMap[storageRecord.ContentID] = storageRecord
	}

	// Convert the content models into the content responses with tags and storage details.
	return models.ConvertContents(content, tagsMap, storageMap), nil
}

func DeleteContents(scanID uint) error {
	if err := database.DB.Where("scan_id = ?", scanID).Delete(&models.ContentModel{}).Error; err != nil {
		return err
	}

	return nil
}

func FindContentByHash(hash string) (models.ContentModel, error) {
	var content models.ContentModel

	if err := database.DB.Where("hashed_body = ?", hash).First(&content).Error; err != nil {
		return content, err
	}

	return content, nil
}

func IncrementAccessCount(contentID uint) error {
	var content models.ContentModel

	if err := database.DB.First(&content, contentID).Error; err != nil {
		return err
	}

	content.AccessCount++

	if err := database.DB.Save(&content).Error; err != nil {
		return err
	}

	return nil
}

func AddContentToScan(scanID uint, contentID uint) error {
	var scan models.ScanModel

	// Retrieve the scan by ID.
	if err := database.DB.First(&scan, scanID).Error; err != nil {
		return err
	}

	// Add the content to the scan's Contents relationship.
	return database.DB.Model(&scan).Association("Contents").Append(&models.ContentModel{ID: contentID})
}
