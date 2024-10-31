package service

import (
	"github.com/codevault-llc/humblebrag-api/internal/database"
	"github.com/codevault-llc/humblebrag-api/internal/models/entities"
)

func GetLicenseByLicense(license string) (entities.LicenseModel, error) {
	var licenseModel entities.LicenseModel

	if err := database.DB.Where("license = ?", license).First(&licenseModel).Error; err != nil {
		return licenseModel, err
	}

	return licenseModel, nil
}
