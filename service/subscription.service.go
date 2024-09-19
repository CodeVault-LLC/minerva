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
)

func CreateSubscription(subscription *models.SubscriptionModel) error {
	tx := constants.DB.Begin()
	if err := tx.Create(&subscription).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func UpdateSubscription(subscription *models.SubscriptionModel) error {
	tx := constants.DB.Begin()
	if err := tx.Save(&subscription).Error; err != nil {
		tx.Rollback()
		return err
	}

	user := &models.UserModel{}
	if err := tx.Model(&models.UserModel{}).Where("id = ?", subscription.UserID).First(user).Error; err != nil {
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

	// Check if the user already has a Stripe customer ID
	if user.StripeCustomerID != "" {
		return user.StripeCustomerID, nil
	}

	// Create a new customer in Stripe
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

	// Update the user record with the new Stripe customer ID
	user.StripeCustomerID = newCustomer.ID
	tx := constants.DB.Begin()
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		return "", err
	}

	return newCustomer.ID, tx.Commit().Error
}

func HandleCheckoutSessionCompleted(checkoutSession *stripe.CheckoutSession) error {
	// Fetch the user from the Stripe customer ID
	user, err := GetUserByStripeCustomerID(checkoutSession.Customer.ID)
	if err != nil {
		return fmt.Errorf("error retrieving user: %v", err)
	}

	// Fetch the subscription from Stripe
	sub, err := subscription.Get(checkoutSession.Subscription.ID, nil)
	if err != nil {
		return fmt.Errorf("error retrieving subscription details from Stripe: %v", err)
	}

	plan := sub.Items.Data[0].Price
	prod, err := product.Get(plan.Product.ID, nil)
	if err != nil {
		return fmt.Errorf("error retrieving product details from Stripe: %v", err)
	}

	// Handle an existing active subscription
	existingSubs, err := GetActiveSubscriptionForUser(user.ID)
	if err == nil && existingSubs != nil && existingSubs.ID != 0 {
		err := CancelExistingSubscription(existingSubs)
		if err != nil {
			return fmt.Errorf("error canceling existing subscription: %v", err)
		}
	}

	// Create a new subscription record
	newSubscription := &models.SubscriptionModel{
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
		return fmt.Errorf("error creating subscription: %v", err)
	}

	notification := &models.NotificationModel{
		Type:    models.NotificationSubscription,
		UserID:  user.ID,
		Message: fmt.Sprintf("You have successfully subscribed to the %s plan.", prod.Name),
		IsRead:  false,
	}

	_, err = CreateNotification(*notification)
	if err != nil {
		return fmt.Errorf("error creating notification: %v", err)
	}

	return nil
}

// GetActiveSubscriptionForUser retrieves the active subscription for a user.
func GetActiveSubscriptionForUser(userID uint) (*models.SubscriptionModel, error) {
	var subscription models.SubscriptionModel

	err := constants.DB.Where("user_id = ? AND status = ?", userID, "active").
		First(&subscription).Attrs(&subscription).Error

	if err != nil {
		return nil, err
	}

	return &subscription, nil
}

// CancelExistingSubscription cancels the current active subscription before creating a new one.
func CancelExistingSubscription(userSubscription *models.SubscriptionModel) error {
	if userSubscription.StripeSubscriptionID == "" {
		return fmt.Errorf("invalid subscription ID")
	}

	params := &stripe.SubscriptionParams{CancelAtPeriodEnd: stripe.Bool(true)}

	_, err := subscription.Update(userSubscription.StripeSubscriptionID, params)
	if err != nil {
		return fmt.Errorf("error updating subscription in Stripe: %v", err)
	}

	userSubscription.Status = "canceled"
	userSubscription.CancelAtPeriodEnd = true

	notification := &models.NotificationModel{
		Type:    models.NotificationSubscription,
		UserID:  userSubscription.UserID,
		Message: "Your current subscription has been canceled.",
		IsRead:  false,
	}

	_, err = CreateNotification(*notification)
	if err != nil {
		return fmt.Errorf("error creating notification: %v", err)
	}

	return UpdateSubscription(userSubscription)
}

func HandleSubscriptionUpdated(sub *stripe.Subscription) error {
	log.Printf("Received subscription update for StripeSubscriptionID: %s", sub.ID)

	existingSubscription, err := retryUntilSubscriptionExists(sub.ID, 5, time.Second)
	if err != nil {
		log.Printf("Error retrieving subscription for StripeSubscriptionID: %s, error: %v", sub.ID, err)
		return err
	}

	if existingSubscription.ID == 0 {
		log.Printf("No subscription found in the database for StripeSubscriptionID: %s", sub.ID)
		return fmt.Errorf("subscription not found or invalid subscription ID")
	}

	log.Printf("Found existing subscription in the database: %v", existingSubscription)

	plan := sub.Items.Data[0].Price
	prod, err := product.Get(plan.Product.ID, nil)
	if err != nil {
		log.Printf("Error retrieving product details from Stripe: %v", err)
		return err
	}

	// Update subscription details
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
	existingSubscription.StripeCustomerID = sub.Customer.ID

	err = UpdateSubscription(existingSubscription)
	if err != nil {
		log.Printf("Error updating subscription in the database: %v", err)
		return err
	}

	log.Printf("Subscription successfully updated in the database: %v", existingSubscription)

	return nil
}

func retryUntilSubscriptionExists(subID string, maxAttempts int, delay time.Duration) (*models.SubscriptionModel, error) {
	for i := 0; i < maxAttempts; i++ {
		subscription, err := GetSubscriptionByStripeSubscriptionID(subID)
		if err == nil && subscription.ID != 0 {
			return &subscription, nil
		}

		log.Printf("Subscription not found, retrying attempt %d/%d", i+1, maxAttempts)
		time.Sleep(delay)
	}

	return nil, fmt.Errorf("subscription not found after %d attempts", maxAttempts)
}
