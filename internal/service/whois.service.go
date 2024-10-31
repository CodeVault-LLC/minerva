package service

import (
	"github.com/codevault-llc/humblebrag-api/internal/database"
	"github.com/codevault-llc/humblebrag-api/internal/models/entities"
)

func CreateWhois(whois entities.WhoisModel) (entities.WhoisModel, error) {
	tx := database.DB.Begin()
	if err := tx.Create(&whois).Error; err != nil {
		tx.Rollback()
		return whois, err
	}

	tx.Commit()
	return whois, nil
}
