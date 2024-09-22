package user

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/codevault-llc/humblebrag-api/helper"
	"github.com/codevault-llc/humblebrag-api/internal/service"
	"github.com/codevault-llc/humblebrag-api/models"
	"github.com/gorilla/mux"
)

func RegisterProfileRoutes(router *mux.Router) {
	router.HandleFunc("/users/me", getCurrentUserHandler).Methods("GET")
	router.HandleFunc("/users/logout", logoutHandler).Methods("POST")
}

// @Summary Get current user
// @Description Get the current user
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} types.Error
// @Failure 404 {object} types.Error
// @Router /users/me [get]
func getCurrentUserHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(models.UserModel)
	if !ok {
		helper.RespondWithError(w, http.StatusUnauthorized, "User not found in context")
		return
	}

	subscription, err := service.GetActiveSubscriptionForUser(user.ID)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve subscriptions")
		return
	}

	notifications, err := service.GetUnreadNotificationsByUserID(user.ID)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve notifications")
		return
	}

	userResponse := models.ConvertUser(user)
	userResponse.Subscription = models.ConvertSubscription(*subscription)
	userResponse.Notifications = models.ConvertNotifications(notifications)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(userResponse); err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to encode user response")
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(models.UserModel)
	if user.ID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	token := r.Header.Get("Authorization")
	if token == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	token = token[7:]

	if err := service.RemoveUserToken(token); err != nil {
		log.Printf("Error removing user token: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
