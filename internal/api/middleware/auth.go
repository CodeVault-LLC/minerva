package middleware

import (
	"log"
	"net/http"

	"github.com/codevault-llc/humblebrag-api/helper"
	"github.com/codevault-llc/humblebrag-api/internal/database/cache"
	"github.com/codevault-llc/humblebrag-api/internal/service"
	"github.com/codevault-llc/humblebrag-api/models"
)

func UserAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if (r.URL.Path == "/api/stripe" && r.Method == "POST") || (r.URL.Path == "/api/auth/discord/callback" && r.Method == "GET") || (r.URL.Path == "/api/auth/discord" && r.Method == "GET") {
			next.ServeHTTP(w, r)
			return
		}

		if cache.SessionManager.Get(r.Context(), "user") == nil {
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
			user, ok := cache.SessionManager.Get(r.Context(), "user").(models.UserModel)
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