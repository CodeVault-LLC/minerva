package repository

import (
	"github.com/codevault-llc/humblebrag-api/internal/models/entities"
	"gorm.io/gorm"
)

type ScanRepository struct {
	db *gorm.DB
}

// NewScanRepository creates a new ScanRepository
func NewScanRepository(db *gorm.DB) *ScanRepository {
	return &ScanRepository{
		db: db,
	}
}

// SaveScanResult saves the scan result in the database
func (repository *ScanRepository) SaveScanResult(job *entities.JobModel, scan entities.ScanModel) error {
	tx := repository.db.Begin()
	if err := tx.Create(&scan).Error; err != nil {
		tx.Rollback()
		return err
	}

	job.ScanID = scan.ID
	if err := tx.Save(&job).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
