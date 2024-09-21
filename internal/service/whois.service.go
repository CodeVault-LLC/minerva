package service

import (
	"github.com/codevault-llc/humblebrag-api/internal/database"
	"github.com/codevault-llc/humblebrag-api/models"
)

func CreateWhois(whois models.WhoisModel) (models.WhoisModel, error) {
	if err := database.DB.Create(&whois).Error; err != nil {
		return whois, err
	}

	return whois, nil
}
