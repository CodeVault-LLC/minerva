package service

import (
	"github.com/codevault-llc/humblebrag-api/internal/database"
	"github.com/codevault-llc/humblebrag-api/models"
)

// GetUserByStripeCustomerID retrieves a user by their Stripe customer ID.
func GetUserByStripeCustomerID(stripeCustomerID string) (*models.UserModel, error) {
	var user models.UserModel

	err := database.DB.Where("stripe_customer_id = ?", stripeCustomerID).
		First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetSubscriptionByStripeSubscriptionID(subscriptionID string) (models.SubscriptionModel, error) {
	var subscription models.SubscriptionModel

	if err := database.DB.Where("stripe_subscription_id = ?", subscriptionID).
		First(&subscription).
		Error; err != nil {
		return subscription, err
	}

	return subscription, nil
}
