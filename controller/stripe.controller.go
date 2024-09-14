package controller

import (
	"github.com/codevault-llc/humblebrag-api/constants"
	"github.com/codevault-llc/humblebrag-api/models"
)

// GetUserByStripeCustomerID retrieves a user by their Stripe customer ID.
func GetUserByStripeCustomerID(stripeCustomerID string) (*models.User, error) {
	var user models.User

	err := constants.DB.Where("stripe_customer_id = ?", stripeCustomerID).
		First(&user).Error

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetSubscriptionByStripeSubscriptionID(subscriptionID string) (models.Subscription, error) {
	var subscription models.Subscription

	if err := constants.DB.Where("stripe_subscription_id = ?", subscriptionID).
		First(&subscription).
		Error; err != nil {
		return subscription, err
	}

	return subscription, nil
}
