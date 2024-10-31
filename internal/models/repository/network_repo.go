package repository

import "gorm.io/gorm"

type NetworkRepository struct {
	db *gorm.DB
}

// NewNetworkRepository creates a new NetworkRepository
func NewNetworkRepository(db *gorm.DB) *NetworkRepository {
	return &NetworkRepository{
		db: db,
	}
}
