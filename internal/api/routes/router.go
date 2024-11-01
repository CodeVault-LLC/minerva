package routes

import (
	"net/http"

	"github.com/codevault-llc/humblebrag-api/internal/api/handlers/scan"
	"github.com/codevault-llc/humblebrag-api/internal/api/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
)

func SetupRouter(app *fiber.App) *fiber.App {
	app.Use(filesystem.New(filesystem.Config{
		Root: http.Dir("./swagger"),
	}))

	app.Get("/docs", serveReDoc)

	api := app.Group("/api/v1")

	// Middlewares
	api.Use(middleware.SubscriptionAuthMiddleware)

	// Scan routes
	_ = scan.RegisterModulesRoutes(api)
	_ = scan.RegisterScanRoutes(api)
	_ = scan.RegisterJobRoutes(api)

	return app
}
