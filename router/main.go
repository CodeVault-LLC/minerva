package router

import (
	"context"
	"encoding/gob"
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
	"github.com/rs/cors"
)

func Start() {
	gob.Register(models.User{})

	router := mux.NewRouter()

	router.HandleFunc("/auth/discord", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, constants.DiscordConfig.AuthCodeURL("random"), http.StatusTemporaryRedirect)
	})

	router.HandleFunc("/auth/discord/callback", func(w http.ResponseWriter, r *http.Request) {
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

		// Check if we have the user, if not then create it
		user := models.User{}
		constants.DB.Where("discord_id = ?", discordUser.Id).First(&user)

		if user.ID == 0 {
			user = models.User{
				DiscordId:    discordUser.Id,
				Username:     discordUser.Username,
				Email:        discordUser.Email,
				Avatar:       discordUser.Avatar,
				AccessToken:  token.AccessToken,
				Provider:     "discord",
				History:      []models.History{},
				Subscriptions: []models.Subscription{},
				Scans:        []models.Scan{},
			}

			user, err = controller.CreateUser(user)
			if err != nil {
				utils.RespondWithError(w, 500, "Error creating user")
				return
			}
		}

		constants.SessionManager.Put(r.Context(), "user", user)
		w.Write([]byte("<script>window.close()</script>"))
	})

	router.Use(func(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if (r.URL.Path == "/auth/discord/callback" && r.Method == "GET") || (r.URL.Path == "/auth/discord" && r.Method == "GET") {
            next.ServeHTTP(w, r)
            return
        }

        user, ok := constants.SessionManager.Get(r.Context(), "user").(models.User)
        if !ok || user.ID == 0 {
						log.Println("Unauthorized")
            utils.RespondWithError(w, 401, "Unauthorized")
            return
        }

        next.ServeHTTP(w, r)
    })
})

	GlobalRouter(router)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	handler := c.Handler(router)

	fmt.Println("Server started on port 3000")

	http.ListenAndServe(":3000", constants.SessionManager.LoadAndSave(handler))
}
