package middleware

import (
	"net/http"

	"github.com/codevault-llc/humblebrag-api/helper"
	"github.com/codevault-llc/humblebrag-api/internal/service"
	"github.com/codevault-llc/humblebrag-api/pkg/responder"
)

func SubscriptionAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("license")
		if token == "" {
			responder.WriteJSONResponse(w, responder.CreateError(responder.ErrAuthInvalidToken))
			return
		}

		license, err := service.GetLicenseByLicense(token)
		if err != nil {
			responder.WriteJSONResponse(w, responder.CreateError(responder.ErrAuthInvalidToken))
			return
		}

		if license.ID == 0 {
			responder.WriteJSONResponse(w, responder.CreateError(responder.ErrAuthInvalidToken))
			return
		}

		r = r.WithContext(helper.AddLicenseToContext(r.Context(), license))
		next.ServeHTTP(w, r)
	})
}
