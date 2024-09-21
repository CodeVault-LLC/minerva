package user

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/codevault-llc/humblebrag-api/internal/service"
	"github.com/codevault-llc/humblebrag-api/models"
	"github.com/codevault-llc/humblebrag-api/pkg/logger"
	"github.com/codevault-llc/humblebrag-api/pkg/utils"
	"github.com/gorilla/mux"
)

func RegisterCheckoutRoutes(router *mux.Router) {
	router.HandleFunc("/users/create-checkout-session", createCheckoutSessionHandler).Methods("POST")
}

func createCheckoutSessionHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(models.UserModel)
	if user.ID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	type PriceRequest struct {
		PriceID string `json:"priceId"`
	}

	var priceRequest PriceRequest
	err := json.NewDecoder(r.Body).Decode(&priceRequest)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if priceRequest.PriceID == "" {
		http.Error(w, "Missing price_id", http.StatusBadRequest)
		return
	}

	stripeCustomerID, err := service.GetOrCreateStripeCustomer(user.ID)
	if err != nil {
		log.Printf("Error getting/creating Stripe customer: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	session, err := utils.CreateCheckoutSession(priceRequest.PriceID, stripeCustomerID)
	if err != nil {
		log.Printf("Error creating checkout session: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(map[string]string{"id": session.ID})
	if err != nil {
		logger.Log.Error("Failed to encode response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
