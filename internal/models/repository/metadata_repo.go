package repository

import (
	"github.com/codevault-llc/humblebrag-api/internal/models/entities"
	"gorm.io/gorm"
)

type MetadataRepo struct {
	db *gorm.DB
}

// NewMetadataRepository creates a new MetadataRepository
func NewMetadataRepository(db *gorm.DB) *MetadataRepo {
	return &MetadataRepo{
		db: db,
	}
}

var MetadataRepository *MetadataRepo

// Create creates a new metadata record in the database
func (repository *MetadataRepo) Create(metadata entities.MetadataModel) (entities.MetadataModel, error) {
	tx := repository.db.Begin()
	if err := tx.Create(&metadata).Error; err != nil {
		tx.Rollback()
		return entities.MetadataModel{}, err
	}

	tx.Commit()
	return metadata, nil
}

func (repository *MetadataRepo) GetMetadataByScanID(scanID uint) (entities.MetadataModel, error) {
	var metadata entities.MetadataModel

	if err := repository.db.Where("scan_id = ?", scanID).First(&metadata).Error; err != nil {
		return entities.MetadataModel{}, err
	}

	return metadata, nil
}
