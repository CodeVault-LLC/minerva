package routes

import (
	"net/http"

	"github.com/codevault-llc/humblebrag-api/internal/contents"
	"github.com/codevault-llc/humblebrag-api/internal/core"
	"github.com/codevault-llc/humblebrag-api/internal/network"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
)

func SetupRouter(app *fiber.App) *fiber.App {
	app.Use(filesystem.New(filesystem.Config{
		Root: http.Dir("./swagger"),
	}))

	app.Get("/docs", serveReDoc)

	api := app.Group("/api/v1")

	// Scan routes
	_ = core.RegisterCoreRouter(api)
	_ = network.RegisterNetworkRouter(api)
	_ = contents.RegisterContentRoutes(api)

	return app
}
