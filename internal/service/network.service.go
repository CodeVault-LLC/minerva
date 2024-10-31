package service

import (
	"github.com/codevault-llc/humblebrag-api/internal/database"
	"github.com/codevault-llc/humblebrag-api/internal/models/entities"
	"github.com/codevault-llc/humblebrag-api/internal/models/viewmodels"
)

// CreateNetwork creates network in the database
func CreateNetwork(network entities.NetworkModel) (entities.NetworkModel, error) {
	if err := database.DB.Create(&network).Error; err != nil {
		return network, err
	}

	return network, nil
}

// GetScanNetwork retrieves network from the database
func GetScanNetwork(scanID uint) (viewmodels.Network, error) {
	var network entities.NetworkModel

	if err := database.DB.Where("scan_id = ?", scanID).
		Preload("Certificates").
		Preload("DNS").
		Preload("Whois").
		First(&network).Error; err != nil {
		return viewmodels.Network{}, err
	}

	return viewmodels.ConvertNetwork(network), nil
}

func UpdateNetwork(network entities.NetworkModel) (entities.NetworkModel, error) {
	if err := database.DB.Save(&network).Error; err != nil {
		return network, err
	}

	return network, nil
}
