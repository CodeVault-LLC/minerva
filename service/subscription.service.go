package service

import (
	"fmt"
	"log"
	"time"

	"github.com/codevault-llc/humblebrag-api/constants"
	"github.com/codevault-llc/humblebrag-api/models"
	"github.com/stripe/stripe-go/v79"
	customer "github.com/stripe/stripe-go/v79/customer"
	"github.com/stripe/stripe-go/v79/product"
	"github.com/stripe/stripe-go/v79/subscription"
	"gorm.io/gorm"
)

func CreateSubscription(subscription *models.Subscription) error {
	tx := constants.DB.Begin()
	if err := tx.Create(&subscription).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func UpdateSubscription(subscription *models.Subscription) error {
	tx := constants.DB.Begin()
	if err := tx.Save(&subscription).Error; err != nil {
		tx.Rollback()
		return err
	}

	user := &models.User{}
	if err := tx.Model(&models.User{}).Where("id = ?", subscription.UserID).First(user).Error; err != nil {
		tx.Rollback()
		return err
	}

	user.StripeCustomerID = subscription.StripeCustomerID
	if err := tx.Save(user).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
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
			"user_id": fmt.Sprintf("%d", user.ID),
		},
	}

	newCustomer, err := customer.New(params)
	if err != nil {
		return "", err
	}

	tx := constants.DB.Begin()
	user.StripeCustomerID = newCustomer.ID
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		return "", err
	}

	return newCustomer.ID, tx.Commit().Error
}

func HandleCheckoutSessionCompleted(checkoutSession *stripe.CheckoutSession) {
	// Fetch the user from the Stripe customer ID
	user, err := GetUserByStripeCustomerID(checkoutSession.Customer.ID)
	if err != nil {
		log.Printf("Error retrieving user: %v", err)
		return
	}

	// Fetch the subscription from Stripe
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

	// Handle an existing active subscription
	existingSubs, err := GetActiveSubscriptionForUser(user.ID)
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Printf("Error retrieving active subscription: %v", err)
		return
	}

	if existingSubs != nil {
		err := CancelExistingSubscription(existingSubs)
		if err != nil {
			log.Printf("Error canceling existing subscription: %v", err)
			return
		}
	}

	// Create a new subscription record
	newSubscription := &models.Subscription{
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

	err = CreateSubscription(newSubscription)
	if err != nil {
		log.Printf("Error creating subscription: %v", err)
	}
}

// GetActiveSubscriptionForUser retrieves the active subscription for a user.
func GetActiveSubscriptionForUser(userID uint) (*models.Subscription, error) {
	var subscription models.Subscription

	err := constants.DB.Where("user_id = ? AND status = ?", userID, "active").
		First(&subscription).Attrs(&subscription).Error

	if err != nil {
		return nil, err
	}

	return &subscription, nil
}

// CancelExistingSubscription cancels the current active subscription before creating a new one.
func CancelExistingSubscription(userSubscription *models.Subscription) error {
	params := &stripe.SubscriptionParams{CancelAtPeriodEnd: stripe.Bool(true)}

	_, err := subscription.Update(userSubscription.StripeSubscriptionID, params)
	if err != nil {
		return err
	}

	userSubscription.Status = "canceled"
	userSubscription.CancelAtPeriodEnd = true

	return UpdateSubscription(userSubscription)
}

// HandleSubscriptionUpdated safely updates the subscription in the database.
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

	err = UpdateSubscription(&existingSubscription)
	if err != nil {
		log.Printf("Error updating subscription: %v", err)
	}
}
