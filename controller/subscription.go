package controller

import (
	"github.com/codevault-llc/humblebrag-api/constants"
	"github.com/codevault-llc/humblebrag-api/models"
)

func CreateSubscription(subscription models.Subscription) (models.Subscription, error) {
	if err := constants.DB.Create(&subscription).Error; err != nil {
		return subscription, err
	}

	return subscription, nil
}

func UpdateSubscription(subscription models.Subscription) (models.Subscription, error) {
	if err := constants.DB.Save(&subscription).Error; err != nil {
		return subscription, err
	}

	return subscription, nil
}
