package models

import (
	"time"

	"gorm.io/gorm"
)

type Subscription struct {
	gorm.Model

	StripeSubscriptionID string `gorm:"unique;not null"`
	StripePriceID        string `gorm:"not null"`
	StripeCustomerID     string `gorm:"not null"`

	UserID uint
	User   User

	PlanType           string  `gorm:"not null"`
	PlanName           string  `gorm:"not null"`
	Price              float64 `gorm:"not null"`
	Currency           string  `gorm:"not null"`
	Interval           string  `gorm:"not null"` // 'month', 'year', etc.
	Status             string  `gorm:"not null"` // 'active', 'canceled', etc.
	CurrentPeriodStart time.Time
	CurrentPeriodEnd   time.Time
	CancelAtPeriodEnd  bool

	History []History `gorm:"foreignKey:SubscriptionID"`
}

type SubscriptionResponse struct {
	ID                 uint      `json:"id"`
	PlanName           string    `json:"plan_name"`
	Price              float64   `json:"price"`
	Currency           string    `json:"currency"`
	Interval           string    `json:"interval"`
	Status             string    `json:"status"`
	CurrentPeriodStart time.Time `json:"current_period_start"`
	CurrentPeriodEnd   time.Time `json:"current_period_end"`
	CancelAtPeriodEnd  bool      `json:"cancel_at_period_end"`
}

func ConvertSubscription(subscription Subscription) SubscriptionResponse {
	return SubscriptionResponse{
		ID:                 subscription.ID,
		PlanName:           subscription.PlanName,
		Price:              subscription.Price,
		Currency:           subscription.Currency,
		Interval:           subscription.Interval,
		Status:             subscription.Status,
		CurrentPeriodStart: subscription.CurrentPeriodStart,
		CurrentPeriodEnd:   subscription.CurrentPeriodEnd,
		CancelAtPeriodEnd:  subscription.CancelAtPeriodEnd,
	}
}
