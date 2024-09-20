package webhook

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/codevault-llc/humblebrag-api/helper"
	"github.com/codevault-llc/humblebrag-api/internal/service"
	"github.com/gorilla/mux"
	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/webhook"
)

func RegisterStripeRoutes(api *mux.Router) {
	api.HandleFunc("/stripe", stripeWebhook).Methods("POST")
}

func stripeWebhook(w http.ResponseWriter, r *http.Request) {
	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)

	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v", err)
		helper.RespondWithError(w, http.StatusServiceUnavailable, "Error reading request body")
		return
	}

	event := stripe.Event{}

	if err := json.Unmarshal(payload, &event); err != nil {
		log.Printf("Error parsing webhook JSON: %v", err)
		helper.RespondWithError(w, http.StatusBadRequest, "Error parsing webhook JSON")
		return
	}

	signatureHeader := r.Header.Get("Stripe-Signature")
	event, err = webhook.ConstructEvent(payload, signatureHeader, os.Getenv("STRIPE_WEBHOOK_SECRET"))

	if err != nil {
		log.Printf("Error verifying webhook signature: %v", err)
		helper.RespondWithError(w, http.StatusBadRequest, "Error verifying webhook signature")
		return
	}

	switch event.Type {
	case "checkout.session.completed":
		var checkoutSession stripe.CheckoutSession
		err := json.Unmarshal(event.Data.Raw, &checkoutSession)
		if err != nil {
			log.Printf("Error parsing checkout session: %v", err)
			helper.RespondWithError(w, http.StatusBadRequest, "Error parsing checkout session")
			return
		}
		err = service.HandleCheckoutSessionCompleted(&checkoutSession)
		if err != nil {
			log.Printf("Error handling checkout session: %v", err)
			helper.RespondWithError(w, http.StatusInternalServerError, "Error handling checkout session")
			return
		}
	case "customer.subscription.updated", "customer.subscription.deleted":
		var sub stripe.Subscription
		err := json.Unmarshal(event.Data.Raw, &sub)
		if err != nil {
			log.Printf("Error parsing subscription: %v", err)
			helper.RespondWithError(w, http.StatusBadRequest, "Error parsing subscription")
			return
		}
		err = service.HandleSubscriptionUpdated(&sub)
		if err != nil {
			log.Printf("Error handling subscription update: %v", err)
			helper.RespondWithError(w, http.StatusInternalServerError, "Error handling subscription update")
			return
		}
	default:
		log.Printf("Unhandled event type: %s", event.Type)
	}

	helper.RespondWithJSON(w, http.StatusOK, map[string]string{"status": "success"})
}
