package controller

import (
	"github.com/codevault-llc/humblebrag-api/constants"
	"github.com/codevault-llc/humblebrag-api/models"
)

func CreateScan(scan models.Scan) (models.Scan, error) {
	if err := constants.DB.Create(&scan).Error; err != nil {
		return scan, err
	}

	return scan, nil
}
