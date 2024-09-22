package service

import (
	"github.com/codevault-llc/humblebrag-api/internal/database"
	"github.com/codevault-llc/humblebrag-api/models"
)

// CreateNetwork creates network in the database
func CreateNetwork(network models.NetworkModel) (models.NetworkModel, error) {
	if err := database.DB.Create(&network).Error; err != nil {
		return network, err
	}

	return network, nil
}

// GetScanNetwork retrieves network from the database
func GetScanNetwork(scanID string) (models.NetworkResponse, error) {
	var network models.NetworkModel

	if err := database.DB.Where("scan_id = ?", scanID).
		Preload("Certificates").
		Preload("Certificates.CertificateResult").
		Preload("Whois").
		First(&network).Error; err != nil {
		return models.NetworkResponse{}, err
	}

	return models.ConvertNetwork(network), nil
}

func UpdateNetwork(network models.NetworkModel) (models.NetworkModel, error) {
	if err := database.DB.Save(&network).Error; err != nil {
		return network, err
	}

	return network, nil
}
