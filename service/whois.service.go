package service

import (
	"github.com/codevault-llc/humblebrag-api/constants"
	"github.com/codevault-llc/humblebrag-api/models"
)

func CreateWhois(whois models.WhoisModel) (models.WhoisModel, error) {
	if err := constants.DB.Create(&whois).Error; err != nil {
		return whois, err
	}

	return whois, nil
}
