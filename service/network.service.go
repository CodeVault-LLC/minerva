package service

import (
	"github.com/codevault-llc/humblebrag-api/constants"
	"github.com/codevault-llc/humblebrag-api/models"
)

// CreateNetwork creates network in the database
func CreateNetwork(network models.NetworkModel) (models.NetworkModel, error) {
	if err := constants.DB.Create(&network).Error; err != nil {
		return network, err
	}

	return network, nil
}

// GetScanNetwork retrieves network from the database
func GetScanNetwork(scanID string) (models.NetworkResponse, error) {
	var network models.NetworkModel

	if err := constants.DB.Where("scan_id = ?", scanID).
		Find(&network).
		Error; err != nil {
		return models.ConvertNetwork(network), err
	}

	return models.ConvertNetwork(network), nil
}
