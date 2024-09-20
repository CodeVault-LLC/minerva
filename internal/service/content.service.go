package service

import (
	"github.com/codevault-llc/humblebrag-api/internal/database"
	"github.com/codevault-llc/humblebrag-api/models"
)

// CreateContent creates content in the database
func CreateContent(content models.ContentModel) (models.ContentModel, error) {
	if err := database.DB.Create(&content).Error; err != nil {
		return content, err
	}

	return content, nil
}

// CreateContents gets content from the database
func GetScanContent(scanID string) ([]models.ContentResponse, error) {
	var content []models.ContentModel

	if err := database.DB.Where("scan_id = ?", scanID).
		Find(&content).
		Error; err != nil {
		return models.ConvertContents(content), err
	}

	return models.ConvertContents(content), nil
}
