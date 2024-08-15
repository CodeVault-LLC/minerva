package user

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/codevault-llc/humblebrag-api/constants"
	"github.com/codevault-llc/humblebrag-api/controller"
	"github.com/codevault-llc/humblebrag-api/models"
	"github.com/codevault-llc/humblebrag-api/utils"
	"github.com/gorilla/mux"
)

func UserRouter(router *mux.Router) {
	router.HandleFunc("/users/me", GetCurrentUser).Methods("GET")

	router.HandleFunc("/auth/discord", func(w http.ResponseWriter, r *http.Request) {
		url := constants.DiscordConfig.AuthCodeURL("random")
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}).Methods("GET")

	router.HandleFunc("/auth/discord/extension", func(w http.ResponseWriter, r *http.Request) {
		url := constants.DiscordConfigExtension.AuthCodeURL("random")
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	}).Methods("GET")

	router.HandleFunc("/auth/discord/callback", AuthenticateDiscord).Methods("GET")
	router.HandleFunc("/auth/discord/callback/extension", AuthenticateDiscordExtension).Methods("GET")

}

// Authenticate for popup windows for our Chrome Extension
func AuthenticateDiscordExtension(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Authenticating Discord for Chrome Extension")
	if r.FormValue("state") != "random" {
		utils.RespondWithError(w, 400, "Invalid State")
		return
	}

	token, err := constants.DiscordConfigExtension.Exchange(r.Context(), r.FormValue("code"))
	if err != nil {
		log.Println(err)
		utils.RespondWithError(w, 500, "Error exchanging token")
		return
	}

	res, err := constants.DiscordConfigExtension.Client(context.Background(), token).Get("https://discord.com/api/users/@me")
	if err != nil || res.StatusCode != 200 {
		w.WriteHeader(http.StatusInternalServerError)
		if err != nil {
			w.Write([]byte(err.Error()))
		} else {
			w.Write([]byte(res.Status))
		}
		return
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	type DiscordUser struct {
		Id            string `json:"id"`
		Username      string `json:"username"`
		Discriminator string `json:"discriminator"`
		Avatar        string `json:"avatar"`
		Verified      bool   `json:"verified"`
		Email         string `json:"email"`
		Flags         int    `json:"flags"`
	}

	var discordUser DiscordUser
	err = json.Unmarshal(body, &discordUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	user := models.User{}
	constants.DB.Where("discord_id = ?", discordUser.Id).First(&user)

	if user.ID == 0 {
		user = models.User{
			DiscordId:     discordUser.Id,
			Username:      discordUser.Username,
			Email:         discordUser.Email,
			Avatar:        discordUser.Avatar,
			AccessToken:   token.AccessToken,
			Provider:      "discord",
			History:       []models.History{},
			Subscriptions: []models.Subscription{},
			Scans:         []models.Scan{},
		}

		user, err = controller.CreateUser(user)
		if err != nil {
			utils.RespondWithError(w, 500, "Error creating user")
			return
		}
	}

	constants.SessionManager.Put(r.Context(), "user", user)
	w.Write([]byte("<script>window.close()</script>"))
}

func GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value("user").(models.User)

	response := utils.ConvertUser(user)

	utils.RespondWithJSON(w, 200, response)
}

// Authenticate for default browsers
func AuthenticateDiscord(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Authenticating Discord")
	if r.FormValue("state") != "random" {
		utils.RespondWithError(w, 400, "Invalid State")
		return
	}

	token, err := constants.DiscordConfig.Exchange(r.Context(), r.FormValue("code"))
	if err != nil {
		log.Println(err)
		utils.RespondWithError(w, 500, "Error exchanging token")
		return
	}

	res, err := constants.DiscordConfig.Client(context.Background(), token).Get("https://discord.com/api/users/@me")
	if err != nil || res.StatusCode != 200 {
		w.WriteHeader(http.StatusInternalServerError)
		if err != nil {
			w.Write([]byte(err.Error()))
		} else {
			w.Write([]byte(res.Status))
		}
		return
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	type DiscordUser struct {
		Id            string `json:"id"`
		Username      string `json:"username"`
		Discriminator string `json:"discriminator"`
		Avatar        string `json:"avatar"`
		Verified      bool   `json:"verified"`
		Email         string `json:"email"`
		Flags         int    `json:"flags"`
	}

	var discordUser DiscordUser
	err = json.Unmarshal(body, &discordUser)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	user := models.User{}
	constants.DB.Where("discord_id = ?", discordUser.Id).First(&user)

	if user.ID == 0 {
		user = models.User{
			DiscordId:     discordUser.Id,
			Username:      discordUser.Username,
			Email:         discordUser.Email,
			Avatar:        discordUser.Avatar,
			AccessToken:   token.AccessToken,
			Provider:      "discord",
			History:       []models.History{},
			Subscriptions: []models.Subscription{},
			Scans:         []models.Scan{},
		}

		user, err = controller.CreateUser(user)
		if err != nil {
			utils.RespondWithError(w, 500, "Error creating user")
			return
		}
	}

	userToken, err := utils.GenerateJWT(user)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	tokenModel := models.UserToken{
		Token:  userToken,
		UserID: user.ID,
	}

	if err := constants.DB.Create(&tokenModel).Error; err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to create token")
		return
	}

	http.Redirect(w, r, fmt.Sprintf("http://localhost:5173/auth?token=%s", userToken), http.StatusTemporaryRedirect)
}
