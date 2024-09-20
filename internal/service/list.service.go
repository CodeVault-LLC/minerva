package service

import (
	"github.com/codevault-llc/humblebrag-api/internal/database"
	"github.com/codevault-llc/humblebrag-api/models"
)

func CreateList(list models.ListModel) (models.ListModel, error) {
	if err := database.DB.Create(&list).Error; err != nil {
		return list, err
	}

	return list, nil
}

func CreateLists(lists []models.ListModel) ([]models.ListModel, error) {
	if err := database.DB.Create(&lists).Error; err != nil {
		return lists, err
	}

	return lists, nil
}
