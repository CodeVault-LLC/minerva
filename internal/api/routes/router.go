package routes

import (
	"encoding/gob"
	"net/http"

	"github.com/codevault-llc/humblebrag-api/internal/api/handlers/scan"
	"github.com/codevault-llc/humblebrag-api/internal/api/middleware"
	"github.com/codevault-llc/humblebrag-api/internal/database/models"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *mux.Router {
	gob.Register(models.LicenseModel{})

	r := mux.NewRouter()

	r.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger/", http.FileServer(http.Dir("./docs"))))

	r.HandleFunc("/docs/", serveReDoc).Methods("GET")

	api := r.PathPrefix("/api/v1").Subrouter()

	// Middlewares
	api.Use(middleware.SubscriptionAuthMiddleware)

	// Scan routes
	scan.RegisterModulesRoutes(api)
	scan.RegisterStatisticsRoutes(api)
	scan.RegisterScanRoutes(api)

	return r
}
