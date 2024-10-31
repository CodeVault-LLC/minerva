package service

import (
	"github.com/codevault-llc/humblebrag-api/internal/database"
	"github.com/codevault-llc/humblebrag-api/internal/database/storage"
	"github.com/codevault-llc/humblebrag-api/internal/models/entities"
	"github.com/codevault-llc/humblebrag-api/internal/models/viewmodels"
	"github.com/codevault-llc/humblebrag-api/pkg/logger"
	"go.uber.org/zap"
)

// CreateContent creates content in the database
func CreateContent(content entities.ContentModel) (entities.ContentModel, error) {
	if err := database.DB.Create(&content).Error; err != nil {
		return content, err
	}

	return content, nil
}

func CreateContentTags(tags []entities.ContentTagsModel) error {
	if err := database.DB.Create(&tags).Error; err != nil {
		return err
	}

	return nil
}

func CreateContentStorage(storage entities.ContentStorageModel) error {
	if err := database.DB.Create(&storage).Error; err != nil {
		return err
	}

	return nil
}

// CreateContents gets content from the database
func GetScanContents(scanID uint) ([]viewmodels.Contents, error) {
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

func GetScanContent(scanID uint, contentID uint) (viewmodels.Content, error) {
	var content entities.ContentModel

	if err := database.DB.First(&content, contentID).Error; err != nil {
		return viewmodels.Content{}, err
	}

	// Go inside the s3 bucket and get the contents.
	contentBody, err := storage.DownloadFile("content-bucket", content.HashedBody)
	if err != nil {
		return viewmodels.Content{}, err
	}

	if err := IncrementAccessCount(contentID); err != nil {
		return viewmodels.Content{}, err
	}

	return viewmodels.Content{
		ID:           content.ID,
		FileSize:     content.FileSize,
		FileType:     content.FileType,
		StorageType:  content.StorageType,
		LastAccessed: content.LastAccessedAt,
		AccessCount:  content.AccessCount,
		Body:         string(contentBody),
	}, nil
}

func DeleteContents(scanID uint) error {
	if err := database.DB.Where("scan_id = ?", scanID).Delete(&entities.ContentModel{}).Error; err != nil {
		return err
	}

	return nil
}

func FindContentByHash(hash string) (entities.ContentModel, error) {
	var content entities.ContentModel

	if err := database.DB.Where("hashed_body = ?", hash).First(&content).Error; err != nil {
		return content, err
	}

	return content, nil
}

func IncrementAccessCount(contentID uint) error {
	var content entities.ContentModel

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
	var scan entities.ScanModel

	// Retrieve the scan by ID.
	if err := database.DB.First(&scan, scanID).Error; err != nil {
		return err
	}

	// Add the content to the scan's Contents relationship.
	return database.DB.Model(&scan).Association("Contents").Append(&entities.ContentModel{ID: contentID})
}
