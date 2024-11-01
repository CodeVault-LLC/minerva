package repository

import (
	"github.com/codevault-llc/humblebrag-api/internal/models/entities"
	"gorm.io/gorm"
)

type NetworkRepo struct {
	db *gorm.DB
}

var NetworkRepository *NetworkRepo

// NewNetworkRepository creates a new NetworkRepository
func NewNetworkRepository(db *gorm.DB) *NetworkRepo {
	return &NetworkRepo{
		db: db,
	}
}

// NetworkRepositoryInterface is the interface for the NetworkRepository
func (n *NetworkRepo) Create(network entities.NetworkModel) (entities.NetworkModel, error) {
	tx := n.db.Begin()
	if err := tx.Create(&network).Error; err != nil {
		tx.Rollback()
		return entities.NetworkModel{}, err
	}

	tx.Commit()
	return network, nil
}

func (n *NetworkRepo) GetScanNetwork(id uint) (entities.NetworkModel, error) {
	var network entities.NetworkModel
	if err := n.db.Where("scan_id = ?", id).First(&network).Error; err != nil {
		return network, err
	}

	return network, nil
}
