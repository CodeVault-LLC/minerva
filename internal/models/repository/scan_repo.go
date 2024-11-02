package repository

import (
	"github.com/codevault-llc/humblebrag-api/internal/models/entities"
	"gorm.io/gorm"
)

type ScanRepo struct {
	db *gorm.DB
}

// NewScanRepository creates a new ScanRepository
func NewScanRepository(db *gorm.DB) *ScanRepo {
	return &ScanRepo{
		db: db,
	}
}

var ScanRepository *ScanRepo

// SaveScanResult saves the scan result in the database
func (repository *ScanRepo) SaveScanResult(job *entities.JobModel, scan entities.ScanModel) (entities.ScanModel, error) {
	tx := repository.db.Begin()
	if err := tx.Create(&scan).Error; err != nil {
		tx.Rollback()
		return entities.ScanModel{}, err
	}

	tx.Commit()
	return scan, nil
}

func (repository *ScanRepo) GetScanResult(scanId uint) (entities.ScanModel, error) {
	var scan entities.ScanModel
	if err := repository.db.Preload("Redirects").First(&scan, scanId).Error; err != nil {
		return scan, err
	}

	return scan, nil
}

func (repository *ScanRepo) GetScans() ([]entities.ScanModel, error) {
	var scans []entities.ScanModel
	if err := repository.db.Find(&scans).Error; err != nil {
		return scans, err
	}

	return scans, nil
}

func (repository *ScanRepo) CompleteScan(scanId uint) error {
	return repository.db.Model(&entities.ScanModel{}).Where("id = ?", scanId).Update("status", entities.ScanStatusComplete).Error
}
