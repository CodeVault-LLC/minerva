package service

import (
	"github.com/codevault-llc/humblebrag-api/internal/database"
	"github.com/codevault-llc/humblebrag-api/internal/database/models"
)

// CreateNmap creates nmap in the database
func CreateNmap(nmap models.NmapModel) (models.NmapModel, error) {
	tx := database.DB.Begin()
	if err := tx.Create(&nmap).Error; err != nil {
		tx.Rollback()
		return nmap, err
	}

	tx.Commit()
	return nmap, nil
}

func GetNmap(scanID string) (models.NmapResponse, error) {
	var nmap models.NmapModel

	if err := database.DB.Where("scan_id = ?", scanID).
		Preload("Hosts").
		First(&nmap).Error; err != nil {
		return models.NmapResponse{}, err
	}

	return models.ConvertNmap(nmap), nil
}
