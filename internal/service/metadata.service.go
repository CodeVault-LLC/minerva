package service

import (
	"github.com/codevault-llc/humblebrag-api/internal/database"
	"github.com/codevault-llc/humblebrag-api/internal/database/models"
)

// CreateMetadata creates metadata in the database
func CreateMetadata(metadata models.MetadataModel) (models.MetadataModel, error) {
	tx := database.DB.Begin()
	if err := tx.Create(&metadata).Error; err != nil {
		tx.Rollback()
		return metadata, err
	}

	tx.Commit()
	return metadata, nil
}

func GetScanMetadataByScanID(scanId uint) (models.MetadataResponse, error) {
	var metadata models.MetadataModel

	if err := database.DB.Where("scan_id = ?", scanId).First(&metadata).Error; err != nil {
		return models.MetadataResponse{}, err
	}

	return models.ConvertMetadata(metadata), nil
}