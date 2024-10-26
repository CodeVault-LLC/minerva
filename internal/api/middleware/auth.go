package middleware

import (
	"net/http"

	"github.com/codevault-llc/humblebrag-api/helper"
	"github.com/codevault-llc/humblebrag-api/internal/service"
)

func SubscriptionAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("license")
		if token == "" {
			helper.RespondWithError(w, 401, "Missing license")
			return
		}

		license, err := service.GetLicenseByLicense(token)
		if err != nil {
			helper.RespondWithError(w, 401, "Invalid license")
			return
		}

		if license.ID == 0 {
			helper.RespondWithError(w, 401, "Invalid license")
			return
		}

		r = r.WithContext(helper.AddLicenseToContext(r.Context(), license))
		next.ServeHTTP(w, r)
	})
}
