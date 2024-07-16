package models

import "gorm.io/gorm"

type History struct {
	gorm.Model

	UserID uint
	User User

	SubscriptionID uint
	Subscription Subscription
}

type HistoryResponse struct {
	ID uint `json:"id"`
	UserID uint `json:"user_id"`
	SubscriptionID uint `json:"subscription_id"`
}
