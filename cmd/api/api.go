package api

import (
	"fmt"
	"net/http"

	"github.com/codevault-llc/humblebrag-api/internal/api/routes"
	"github.com/gorilla/handlers"
	"github.com/rs/cors"
)

// @title Humblebrag-API
// @version 1.0
// @description Humblebrag is a scanner service which detects and reports on the presence of sensitive data in your codebase and infrastructure.
// @termsOfService http://swagger.io/terms/
// @contact.name Humblebrag LLC Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:3000
// @BasePath /api/v1
func Start() {
	api := routes.SetupRouter()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	handler := handlers.CompressHandler(c.Handler(api))

	fmt.Println("Server started on port 3000")
	err := http.ListenAndServe(":3000", handler)
	if err != nil {
		panic(err)
	}
}
