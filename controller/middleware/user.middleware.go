package middleware

import (
	"log"
	"net/http"

	"github.com/codevault-llc/humblebrag-api/constants"
	"github.com/codevault-llc/humblebrag-api/helper"
	"github.com/codevault-llc/humblebrag-api/models"
	"github.com/codevault-llc/humblebrag-api/service"
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
				helper.RespondWithError(w, 401, "Missing token")
				return
			}

			token = token[7:]

			claims, err := helper.ValidateJWT(token)

			if err != nil {
				helper.RespondWithError(w, 401, "Invalid token")
				return
			}

			userToken, err := service.IsValidUserToken(token)
			if err != nil {
				helper.RespondWithError(w, 401, "Invalid token")
				return
			}

			if userToken.ID == 0 {
				helper.RespondWithError(w, 401, "Invalid token")
				return
			}

			userIDFloat, ok := claims["id"].(float64)
			if !ok {
				helper.RespondWithError(w, 401, "Invalid token")
				return
			}

			userID := uint(userIDFloat)

			user, err := service.GetUserById(userID)
			if err != nil {
				helper.RespondWithError(w, 401, "Invalid token")
				return
			}

			if user.ID == 0 {
				helper.RespondWithError(w, 401, "Invalid token")
				return
			}

			r = r.WithContext(helper.AddUserToContext(r.Context(), user))
			next.ServeHTTP(w, r)
		} else {
			user, ok := constants.SessionManager.Get(r.Context(), "user").(models.UserModel)
			if !ok || user.ID == 0 {
				log.Println("Unauthorized")
				helper.RespondWithError(w, 401, "Unauthorized")
				return
			}

			r = r.WithContext(helper.AddUserToContext(r.Context(), user))
			next.ServeHTTP(w, r)
		}
	})
}
