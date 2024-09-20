package user

import (
	"log"
	"net/http"

	"github.com/codevault-llc/humblebrag-api/internal/service"
	"github.com/codevault-llc/humblebrag-api/models"
	"github.com/gorilla/mux"
)

func RegisterSubscriptionRoutes(router *mux.Router) {
	router.HandleFunc("/users/cancel-subscription", cancelSubscriptionHandler).Methods("POST")
}

func cancelSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(models.UserModel)
	if user.ID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	subscriptionFromUser, err := service.GetActiveSubscriptionForUser(user.ID)
	if err != nil {
		log.Printf("Error getting subscriptions: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if subscriptionFromUser.ID == 0 {
		http.Error(w, "No active subscriptions", http.StatusBadRequest)
		return
	}

	err = service.CancelExistingSubscription(subscriptionFromUser)
	if err != nil {
		log.Printf("Error cancelling subscription: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
