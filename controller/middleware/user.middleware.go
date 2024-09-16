package middleware

import (
	"log"
	"net/http"

	"github.com/codevault-llc/humblebrag-api/constants"
	"github.com/codevault-llc/humblebrag-api/models"
	"github.com/codevault-llc/humblebrag-api/service"
	"github.com/codevault-llc/humblebrag-api/utils"
)

func UserAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if (r.URL.Path == "/stripe" && r.Method == "POST") || (r.URL.Path == "/auth/discord/callback" && r.Method == "GET") || (r.URL.Path == "/auth/discord/callback/extension" && r.Method == "GET") || (r.URL.Path == "/auth/discord/extension" && r.Method == "GET") || (r.URL.Path == "/auth/discord" && r.Method == "GET") {
			next.ServeHTTP(w, r)
			return
		}

		if constants.SessionManager.Get(r.Context(), "user") == nil {
			token := r.Header.Get("Authorization")
			if token == "" {
				utils.RespondWithError(w, 401, "Missing token")
				return
			}

			token = token[7:]

			claims, err := utils.ValidateJWT(token)

			if err != nil {
				utils.RespondWithError(w, 401, "Invalid token")
				return
			}

			userToken, err := service.IsValidUserToken(token)
			if err != nil {
				utils.RespondWithError(w, 401, "Invalid token")
				return
			}

			if userToken.ID == 0 {
				utils.RespondWithError(w, 401, "Invalid token")
				return
			}

			user := models.User{}
			constants.DB.First(&user, claims["id"])

			if user.ID == 0 {
				utils.RespondWithError(w, 401, "Invalid token")
				return
			}

			r = r.WithContext(utils.AddUserToContext(r.Context(), user))
			next.ServeHTTP(w, r)
		} else {
			user, ok := constants.SessionManager.Get(r.Context(), "user").(models.User)
			if !ok || user.ID == 0 {
				log.Println("Unauthorized")
				utils.RespondWithError(w, 401, "Unauthorized")
				return
			}

			r = r.WithContext(utils.AddUserToContext(r.Context(), user))
			next.ServeHTTP(w, r)
		}
	})
}
