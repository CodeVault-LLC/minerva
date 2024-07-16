package webhook

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/codevault-llc/humblebrag-api/controller"
	"github.com/codevault-llc/humblebrag-api/models"
	"github.com/codevault-llc/humblebrag-api/utils"
	"github.com/gorilla/mux"
	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/webhook"
)

func WebhookRouter(router *mux.Router) {
	router.HandleFunc("/stripe", StripeWebhook).Methods("POST")
}

func StripeWebhook(w http.ResponseWriter, r *http.Request) {
	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)

	payload, err := ioutil.ReadAll(r.Body)

	if err != nil {
		utils.RespondWithError(w, http.StatusServiceUnavailable, "Request body too large")
		return
	}

	event := stripe.Event{}

	if err := json.Unmarshal(payload, &event); err != nil {
    fmt.Fprintf(os.Stderr, "⚠️  Webhook error while parsing basic request. %v\n", err.Error())
    w.WriteHeader(http.StatusBadRequest)
    return
  }

	signatureHeader := r.Header.Get("Stripe-Signature")
	event, err = webhook.ConstructEvent(payload, signatureHeader, os.Getenv("STRIPE_WEBHOOK_SECRET"))

	if err != nil {
    fmt.Fprintf(os.Stderr, "⚠️  Webhook signature verification failed. %v\n", err)
    w.WriteHeader(http.StatusBadRequest) // Return a 400 error on a bad signature
    return
  }

	switch event.Type {
		case "customer.subscription.created":
			var subscriptionResponse stripe.Subscription

			err := json.Unmarshal(event.Data.Raw, &subscriptionResponse)

			if err != nil {
				utils.RespondWithError(w, http.StatusServiceUnavailable, "Error parsing subscription")
				return
			}

			user, err := controller.GetUserByEmail(subscriptionResponse.Customer.Email)
			if err != nil {
				utils.RespondWithError(w, http.StatusServiceUnavailable, "Error fetching user")
				return
			}

			if user.ID == 0 {
				utils.RespondWithError(w, http.StatusServiceUnavailable, "User not found")
				return
			}

			if user.Subscriptions != nil {
				utils.RespondWithError(w, http.StatusServiceUnavailable, "User already has a subscription")
				return
			}

			subscription, err := controller.CreateSubscription(models.Subscription{
				UserID: user.ID,
				SubscriptionId: subscriptionResponse.ID,
				PlanType: subscriptionResponse.Items.Data[0].Plan.Nickname,
				Price: subscriptionResponse.Items.Data[0].Plan.AmountDecimal,
				// Status of the subscription aka active, canceled, etc
				Status: string(subscriptionResponse.Status),
			})

			if err != nil {
				utils.RespondWithError(w, http.StatusServiceUnavailable, "Error creating subscription")
				return
			}

			utils.RespondWithJSON(w, http.StatusOK, subscription)
			return

		case "customer.subscription.deleted":
			var subscription stripe.Subscription

			err := json.Unmarshal(event.Data.Raw, &subscription)

			if err != nil {
				utils.RespondWithError(w, http.StatusServiceUnavailable, "Error parsing subscription")
				return
			}

			user, err := controller.GetUserByEmail(subscription.Customer.Email)

			if err != nil {
				utils.RespondWithError(w, http.StatusServiceUnavailable, "Error fetching user")
				return
			}

			if user.ID == 0 {
				utils.RespondWithError(w, http.StatusServiceUnavailable, "User not found")
				return
			}

			if user.Subscriptions == nil {
				utils.RespondWithError(w, http.StatusServiceUnavailable, "User does not have a subscription")
				return
			}

			sub, err := controller.UpdateSubscription(models.Subscription{
				UserID: user.ID,
				SubscriptionId: subscription.ID,
				PlanType: subscription.Items.Data[0].Plan.Nickname,
				Price: subscription.Items.Data[0].Plan.AmountDecimal,
				Status: string(subscription.Status),
			})

			if err != nil {
				utils.RespondWithError(w, http.StatusServiceUnavailable, "Error updating subscription")
				return
			}

			utils.RespondWithJSON(w, http.StatusOK, sub)
			return

		case "customer.subscription.updated":
			var subscription stripe.Subscription

			err := json.Unmarshal(event.Data.Raw, &subscription)

			if err != nil {
				utils.RespondWithError(w, http.StatusServiceUnavailable, "Error parsing subscription")
				return
			}

			user, err := controller.GetUserByEmail(subscription.Customer.Email)
			if err != nil {
				utils.RespondWithError(w, http.StatusServiceUnavailable, "Error fetching user")
				return
			}

			if user.ID == 0 {
				utils.RespondWithError(w, http.StatusServiceUnavailable, "User not found")
				return
			}

			if user.Subscriptions == nil {
				utils.RespondWithError(w, http.StatusServiceUnavailable, "User does not have a subscription")
				return
			}

			sub, err := controller.UpdateSubscription(models.Subscription{
				UserID: user.ID,
				SubscriptionId: subscription.ID,
				PlanType: subscription.Items.Data[0].Plan.Nickname,
				Price: subscription.Items.Data[0].Plan.AmountDecimal,
				Status: string(subscription.Status),
			})

			if err != nil {
				utils.RespondWithError(w, http.StatusServiceUnavailable, "Error updating subscription")
				return
			}

			utils.RespondWithJSON(w, http.StatusOK, sub)
			return

		default:
			fmt.Fprintf(os.Stderr, "⚠️  Webhook received unknown event type: %s\n", event.Type)
			w.WriteHeader(http.StatusBadRequest)
			return
	}
}
