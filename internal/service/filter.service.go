package service

import (
	"github.com/codevault-llc/humblebrag-api/internal/database"
	"github.com/codevault-llc/humblebrag-api/internal/database/models"
)

func CreateFilter(list models.FilterModel) (models.FilterModel, error) {
	if err := database.DB.Create(&list).Error; err != nil {
		return list, err
	}

	return list, nil
}

func CreateFilters(lists []models.FilterModel) ([]models.FilterModel, error) {
	if err := database.DB.Create(&lists).Error; err != nil {
		return lists, err
	}

	return lists, nil
}
