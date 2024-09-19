package models

import "gorm.io/gorm"

type HistoryModel struct {
	gorm.Model

	UserID uint
	User   UserModel

	SubscriptionID uint
	Subscription   SubscriptionModel
}

type HistoryResponse struct {
	ID             uint `json:"id"`
	UserID         uint `json:"user_id"`
	SubscriptionID uint `json:"subscription_id"`
}
