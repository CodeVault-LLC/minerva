package controller

import (
	"log"
	"time"

	"github.com/codevault-llc/humblebrag-api/constants"
	"github.com/codevault-llc/humblebrag-api/models"
	"github.com/stripe/stripe-go/v79"
	customer "github.com/stripe/stripe-go/v79/customer"
	"github.com/stripe/stripe-go/v79/product"
	"github.com/stripe/stripe-go/v79/subscription"
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

func GetOrCreateStripeCustomer(userID uint) (string, error) {
	user, err := GetUserById(userID)
	if err != nil {
		return "", err
	}

	if user.StripeCustomerID != "" {
		return user.StripeCustomerID, nil
	}

	params := &stripe.CustomerParams{
		Email: stripe.String(user.Email),
		Metadata: map[string]string{
			"user_id": string(user.ID),
		},
	}

	newCustomer, err := customer.New(params)
	if err != nil {
		return "", err
	}

	user.StripeCustomerID = newCustomer.ID
	if err := constants.DB.Save(&user).Error; err != nil {
		return "", err
	}

	if err != nil {
		return "", err
	}

	return newCustomer.ID, nil
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

func GetUserByStripeCustomerID(customerID string) (models.User, error) {
	var user models.User

	if err := constants.DB.Where("stripe_customer_id = ?", customerID).
		First(&user).
		Error; err != nil {
		return user, err
	}

	return user, nil
}

func HandleCheckoutSessionCompleted(checkoutSession *stripe.CheckoutSession) {
	user, err := GetUserByStripeCustomerID(checkoutSession.Customer.ID)
	if err != nil {
		log.Printf("Error retrieving user: %v", err)
		return
	}

	sub, err := subscription.Get(checkoutSession.Subscription.ID, nil)
	if err != nil {
		log.Printf("Error retrieving subscription details from Stripe: %v", err)
		return
	}

	plan := sub.Items.Data[0].Price

	prod, err := product.Get(plan.Product.ID, nil)
	if err != nil {
		log.Printf("Error retrieving product details from Stripe: %v", err)
		return
	}

	// Create a new subscription record
	newSubscription := models.Subscription{
		StripeSubscriptionID: sub.ID,
		StripePriceID:        plan.ID,
		StripeCustomerID:     checkoutSession.Customer.ID,
		UserID:               user.ID,
		PlanName:             prod.Name,
		PlanType:             plan.Nickname,
		Price:                float64(plan.UnitAmount) / 100,
		Currency:             string(plan.Currency),
		Interval:             string(plan.Recurring.Interval),
		Status:               string(sub.Status),
		CurrentPeriodStart:   time.Unix(sub.CurrentPeriodStart, 0),
		CurrentPeriodEnd:     time.Unix(sub.CurrentPeriodEnd, 0),
		CancelAtPeriodEnd:    sub.CancelAtPeriodEnd,
	}

	_, err = CreateSubscription(newSubscription)
	if err != nil {
		log.Printf("Error creating subscription: %v", err)
		return
	}
}

func HandleSubscriptionUpdated(sub *stripe.Subscription) {
	existingSubscription, err := GetSubscriptionByStripeSubscriptionID(sub.ID)
	if err != nil {
		log.Printf("Error retrieving subscription: %v", err)
		return
	}

	plan := sub.Items.Data[0].Price

	prod, err := product.Get(plan.Product.ID, nil)
	if err != nil {
		log.Printf("Error retrieving product details from Stripe: %v", err)
		return
	}

	existingSubscription.StripePriceID = plan.ID
	existingSubscription.PlanName = prod.Name
	existingSubscription.PlanType = plan.Nickname
	existingSubscription.Price = float64(plan.UnitAmount) / 100
	existingSubscription.Currency = string(plan.Currency)
	existingSubscription.Interval = string(plan.Recurring.Interval)
	existingSubscription.Status = string(sub.Status)
	existingSubscription.CurrentPeriodStart = time.Unix(sub.CurrentPeriodStart, 0)
	existingSubscription.CurrentPeriodEnd = time.Unix(sub.CurrentPeriodEnd, 0)
	existingSubscription.CancelAtPeriodEnd = sub.CancelAtPeriodEnd

	_, err = UpdateSubscription(existingSubscription)
	if err != nil {
		log.Printf("Error updating subscription: %v", err)
		return
	}
}

func GetSubscriptionFromUser(userID uint) (models.Subscription, error) {
	var subscription models.Subscription

	if err := constants.DB.Where("user_id = ?", userID).
		First(&subscription).
		Error; err != nil {
		return subscription, err
	}

	return subscription, nil
}
