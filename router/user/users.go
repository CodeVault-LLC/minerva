package user

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/codevault-llc/humblebrag-api/constants"
	"github.com/codevault-llc/humblebrag-api/controller"
	"github.com/codevault-llc/humblebrag-api/models"
	"github.com/codevault-llc/humblebrag-api/utils"
	"github.com/gorilla/mux"
)

func RegisterUserRoutes(router *mux.Router) {
	router.HandleFunc("/users/me", getCurrentUserHandler).Methods("GET")
	router.HandleFunc("/auth/discord", discordAuthRedirectHandler).Methods("GET")
	router.HandleFunc("/auth/discord/extension", discordExtensionAuthRedirectHandler).Methods("GET")
	router.HandleFunc("/auth/discord/callback", discordAuthCallbackHandler).Methods("GET")
	router.HandleFunc("/auth/discord/callback/extension", discordExtensionCallbackHandler).Methods("GET")
	router.HandleFunc("/users/create-checkout-session", createCheckoutSessionHandler).Methods("POST")
	router.HandleFunc("/users/cancel-subscription", cancelSubscriptionHandler).Methods("POST")
	router.HandleFunc("/users/logout", logoutHandler).Methods("POST")
}

func discordAuthRedirectHandler(w http.ResponseWriter, r *http.Request) {
	url := constants.DiscordConfig.AuthCodeURL("random")
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func discordExtensionAuthRedirectHandler(w http.ResponseWriter, r *http.Request) {
	url := constants.DiscordConfigExtension.AuthCodeURL("random")
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func discordExtensionCallbackHandler(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != "random" {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid state parameter")
		return
	}

	token, err := constants.DiscordConfigExtension.Exchange(r.Context(), r.FormValue("code"))
	if err != nil {
		log.Println("Error exchanging token:", err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to exchange token")
		return
	}

	userInfo, err := controller.FetchDiscordUserInfo(*token)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve user info")
		return
	}

	user, err := controller.FindOrCreateUserFromDiscord(*userInfo, token)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create or find user")
		return
	}

	constants.SessionManager.Put(r.Context(), "user", user)
	w.Write([]byte("<script>window.close()</script>"))
}

func discordAuthCallbackHandler(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != "random" {
		utils.RespondWithError(w, http.StatusBadRequest, "Invalid state parameter")
		return
	}

	token, err := constants.DiscordConfig.Exchange(r.Context(), r.FormValue("code"))
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to exchange token")
		return
	}

	userInfo, err := controller.FetchDiscordUserInfo(*token)
	if err != nil {
		fmt.Println(err)
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve user info")
		return
	}

	user, err := controller.FindOrCreateUserFromDiscord(*userInfo, token)

	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create or find user")
		return
	}

	userToken, err := utils.GenerateJWT(user)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to generate JWT")
		return
	}

	if err := controller.SaveUserToken(userToken, user.ID); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to save user token")
		return
	}

	http.Redirect(w, r, fmt.Sprintf("http://localhost:5173/auth?token=%s", userToken), http.StatusTemporaryRedirect)
}

func createCheckoutSessionHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(models.User)
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

	stripeCustomerID, err := controller.GetOrCreateStripeCustomer(user.ID)
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

	json.NewEncoder(w).Encode(map[string]string{"id": session.ID})
}

func cancelSubscriptionHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(models.User)
	if user.ID == 0 {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	subscriptionFromUser, err := controller.GetActiveSubscriptionForUser(user.ID)
	if err != nil {
		log.Printf("Error getting subscriptions: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if subscriptionFromUser.ID == 0 {
		http.Error(w, "No active subscriptions", http.StatusBadRequest)
		return
	}

	err = controller.CancelExistingSubscription(subscriptionFromUser)
	if err != nil {
		log.Printf("Error cancelling subscription: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func getCurrentUserHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(models.User)
	if !ok {
		utils.RespondWithError(w, http.StatusUnauthorized, "User not found in context")
		return
	}

	subscription, err := controller.GetActiveSubscriptionForUser(user.ID)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve subscriptions")
		return
	}

	userResponse := utils.ConvertUser(user)
	userResponse.Subscription = utils.ConvertSubscription(*subscription)

	if err := json.NewEncoder(w).Encode(userResponse); err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to encode user response")
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(models.User)
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

	if err := controller.RemoveUserToken(token); err != nil {
		log.Printf("Error removing user token: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
