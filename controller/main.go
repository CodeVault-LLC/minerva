package controller

import (
	"encoding/gob"
	"fmt"
	"net/http"

	"github.com/codevault-llc/humblebrag-api/constants"
	"github.com/codevault-llc/humblebrag-api/controller/middleware"
	"github.com/codevault-llc/humblebrag-api/models"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func Start() {
	gob.Register(models.UserModel{})
	router := mux.NewRouter()

	router.Use(middleware.UserAuthMiddleware)
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

func GlobalRouter(router *mux.Router) {
	ScanRouter(router)
	RegisterUserRoutes(router)
	WebhookRouter(router)
	RegisterNotificationRoutes(router)
}
