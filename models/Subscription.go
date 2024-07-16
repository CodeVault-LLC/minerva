package models

import "gorm.io/gorm"

type Subscription struct {
	gorm.Model

	SubscriptionId string `gorm:"unique;not null"`

	UserID uint
	User User

	PlanType string `gorm:"not null"`
	Price float64 `gorm:"not null"`
	Status string `gorm:"not null"`

	History []History `gorm:"foreignKey:SubscriptionID"`
}

type SubscriptionResponse struct {
	ID uint `json:"id"`
	SubscriptionId string `json:"subscription_id"`
	PlanType string `json:"plan_type"`
	Price float64 `json:"price"`
	Status string `json:"status"`
}
