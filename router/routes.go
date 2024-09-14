package router

import (
	"github.com/gorilla/mux"

	"github.com/codevault-llc/humblebrag-api/router/scan"
	"github.com/codevault-llc/humblebrag-api/router/user"
	"github.com/codevault-llc/humblebrag-api/router/webhook"
)

func GlobalRouter(router *mux.Router) {
	scan.ScanRouter(router)
	user.RegisterUserRoutes(router)
	webhook.WebhookRouter(router)
}
