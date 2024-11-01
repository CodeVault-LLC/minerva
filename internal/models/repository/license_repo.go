package repository

import (
	"github.com/codevault-llc/humblebrag-api/internal/models/entities"
	"gorm.io/gorm"
)

type LicenseRepo struct {
	db *gorm.DB
}

var LicenseRepository *LicenseRepo

// NewLicenseRepository creates a new LicenseRepository
func NewLicenseRepository(db *gorm.DB) *LicenseRepo {
	return &LicenseRepo{
		db: db,
	}
}

// FindByID retrieves a license by its ID
func (repository *LicenseRepo) FindByID(id uint) (*entities.LicenseModel, error) {
	var license entities.LicenseModel
	err := repository.db.First(&license, id).Error
	return &license, err
}

func (repository *LicenseRepo) GetLicenseByLicense(license string) (*entities.LicenseModel, error) {
	var licenseModel entities.LicenseModel
	err := repository.db.Where("license = ?", license).First(&licenseModel).Error
	return &licenseModel, err
}
