package api

import (
	"github.com/codevault-llc/humblebrag-api/internal/api/routes"
	"github.com/codevault-llc/humblebrag-api/pkg/responder"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
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
	app := fiber.New(fiber.Config{
		Prefork:       false,
		CaseSensitive: true,
		ErrorHandler:  responder.ErrorHandler,
	})

	app.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed, // 1
	}))
	app.Use(helmet.New())

	app.Use(func(c *fiber.Ctx) error {
		c.Set("Content-Security-Policy", "default-src 'self'")
		c.Set("Strict-Transport-Security", "max-age=31536000")
		c.Set("X-Frame-Options", "DENY")
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-XSS-Protection", "1; mode=block")

		return c.Next()
	})

	app.Use(limiter.New(limiter.Config{
		Max:        60,
		Expiration: 60,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return responder.CreateError(responder.ErrLimitReached).Error
		},
	}))

	app.Use(recover.New())
	app.Use(logger.New())

	app.Use(cors.New(cors.Config{
		AllowCredentials: false,
		MaxAge:           3600,
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders:     "Origin, Content-Type, Accept",
		ExposeHeaders:    "",
	}))

	api := routes.SetupRouter(app)

	err := api.Listen(":3000")
	if err != nil {
		panic(err)
	}
}
