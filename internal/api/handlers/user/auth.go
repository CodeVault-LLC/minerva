package user

import (
	"fmt"
	"net/http"

	"github.com/codevault-llc/humblebrag-api/config"
	"github.com/codevault-llc/humblebrag-api/helper"
	"github.com/codevault-llc/humblebrag-api/internal/service"
	"github.com/gorilla/mux"
)

func RegisterAuthRoutes(router *mux.Router) {
	router.HandleFunc("/auth/discord", discordAuthRedirectHandler).Methods("GET")
	router.HandleFunc("/auth/discord/callback", discordAuthCallbackHandler).Methods("GET")
}

func discordAuthRedirectHandler(w http.ResponseWriter, r *http.Request) {
	url := config.DiscordConfig.AuthCodeURL("random")
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func discordAuthCallbackHandler(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != "random" {
		helper.RespondWithError(w, http.StatusBadRequest, "Invalid state parameter")
		return
	}

	token, err := config.DiscordConfig.Exchange(r.Context(), r.FormValue("code"))
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to exchange token")
		return
	}

	userInfo, err := service.FetchDiscordUserInfo(*token)
	if err != nil {
		fmt.Println(err)
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to retrieve user info")
		return
	}

	user, err := service.FindOrCreateUserFromDiscord(*userInfo, token)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to create or find user")
		return
	}

	userToken, err := helper.GenerateJWT(user.ID)
	if err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to generate JWT")
		return
	}

	if err := service.SaveUserToken(userToken, user.ID); err != nil {
		helper.RespondWithError(w, http.StatusInternalServerError, "Failed to save user token")
		return
	}

	http.Redirect(w, r, fmt.Sprintf("http://localhost:5173/auth?token=%s", userToken), http.StatusTemporaryRedirect)
}
