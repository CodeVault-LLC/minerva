package routes

import (
	"encoding/gob"

	"github.com/codevault-llc/humblebrag-api/internal/api/handlers/notification"
	"github.com/codevault-llc/humblebrag-api/internal/api/handlers/scan"
	"github.com/codevault-llc/humblebrag-api/internal/api/handlers/user"
	"github.com/codevault-llc/humblebrag-api/internal/api/handlers/webhook"
	"github.com/codevault-llc/humblebrag-api/internal/api/middleware"
	"github.com/codevault-llc/humblebrag-api/models"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *mux.Router {
	gob.Register(models.UserModel{})
	r := mux.NewRouter()

	api := r.PathPrefix("/api").Subrouter()

	// Middlewares
	api.Use(middleware.UserAuthMiddleware)

	// User routes
	user.RegisterAuthRoutes(api)
	user.RegisterProfileRoutes(api)
	user.RegisterSubscriptionRoutes(api)
	user.RegisterCheckoutRoutes(api)

	// Notification routes
	notification.RegisterNotificationRoutes(api)

	// Scan routes
	scan.RegisterModulesRoutes(api)
	scan.RegisterScanRoutes(api)
	scan.RegisterStatisticsRoutes(api)

	// Webhook routes
	webhook.RegisterStripeRoutes(api)

	return r
}
