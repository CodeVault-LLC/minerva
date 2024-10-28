package service

import (
	"github.com/codevault-llc/humblebrag-api/internal/database"
	"github.com/codevault-llc/humblebrag-api/internal/database/models"
)

func GetLicenseByLicense(license string) (models.LicenseModel, error) {
	var licenseModel models.LicenseModel

	if err := database.DB.Where("license = ?", license).First(&licenseModel).Error; err != nil {
		return licenseModel, err
	}

	return licenseModel, nil
}
