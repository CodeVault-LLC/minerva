package repository

import (
	"github.com/codevault-llc/humblebrag-api/internal/models/entities"
	"gorm.io/gorm"
)

type LicenseRepository struct {
	db *gorm.DB
}

// NewLicenseRepository creates a new LicenseRepository
func NewLicenseRepository(db *gorm.DB) *LicenseRepository {
	return &LicenseRepository{
		db: db,
	}
}

// FindByID retrieves a license by its ID
func (repository *LicenseRepository) FindByID(id uint) (*entities.LicenseModel, error) {
	var license entities.LicenseModel
	err := repository.db.First(&license, id).Error
	return &license, err
}
