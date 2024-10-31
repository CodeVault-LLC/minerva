package service

import (
	"github.com/codevault-llc/humblebrag-api/internal/database"
	"github.com/codevault-llc/humblebrag-api/internal/models/entities"
)

// CreateDNS creates DNS in the database
func CreateDNS(dns entities.DNSModel) (entities.DNSModel, error) {
	tx := database.DB.Begin()
	if err := tx.Create(&dns).Error; err != nil {
		tx.Rollback()
		return dns, err
	}

	tx.Commit()
	return dns, nil
}
