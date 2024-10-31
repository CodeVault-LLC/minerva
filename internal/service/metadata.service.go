package service

import (
	"github.com/codevault-llc/humblebrag-api/internal/database"
	"github.com/codevault-llc/humblebrag-api/internal/models/entities"
	"github.com/codevault-llc/humblebrag-api/internal/models/viewmodels"
)

// CreateMetadata creates metadata in the database
func CreateMetadata(metadata entities.MetadataModel) (entities.MetadataModel, error) {
	tx := database.DB.Begin()
	if err := tx.Create(&metadata).Error; err != nil {
		tx.Rollback()
		return metadata, err
	}

	tx.Commit()
	return metadata, nil
}

func GetScanMetadataByScanID(scanId uint) (viewmodels.Metadata, error) {
	var metadata entities.MetadataModel

	if err := database.DB.Where("scan_id = ?", scanId).First(&metadata).Error; err != nil {
		return viewmodels.Metadata{}, err
	}

	return viewmodels.ConvertMetadata(metadata), nil
}
