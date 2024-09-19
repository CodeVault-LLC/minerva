package service

import (
	"github.com/codevault-llc/humblebrag-api/constants"
	"github.com/codevault-llc/humblebrag-api/models"
)

func CreateList(list models.ListModel) (models.ListModel, error) {
	if err := constants.DB.Create(&list).Error; err != nil {
		return list, err
	}

	return list, nil
}

func CreateLists(lists []models.ListModel) ([]models.ListModel, error) {
	if err := constants.DB.Create(&lists).Error; err != nil {
		return lists, err
	}

	return lists, nil
}
