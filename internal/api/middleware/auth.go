package middleware

import (
	"log"
	"net/http"

	"github.com/codevault-llc/humblebrag-api/helper"
	"github.com/codevault-llc/humblebrag-api/internal/database/cache"
	"github.com/codevault-llc/humblebrag-api/internal/database/models"
	"github.com/codevault-llc/humblebrag-api/internal/service"
)

func SubscriptionAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if cache.SessionManager.Get(r.Context(), "license") == nil {
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
		} else {
			license, ok := cache.SessionManager.Get(r.Context(), "license").(models.LicenseModel)
			if !ok || license.ID == 0 {
				log.Println("Unauthorized")
				helper.RespondWithError(w, 401, "Unauthorized")
				return
			}

			r = r.WithContext(helper.AddLicenseToContext(r.Context(), license))
			next.ServeHTTP(w, r)
		}
	})
}
