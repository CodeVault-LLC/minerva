package middleware

import (
	"github.com/codevault-llc/humblebrag-api/internal/service"
	"github.com/codevault-llc/humblebrag-api/pkg/responder"
	"github.com/gofiber/fiber/v2"
)

func SubscriptionAuthMiddleware(c *fiber.Ctx) error {
	token := c.Get("license")
	if token == "" {
		return responder.CreateError(responder.ErrAuthInvalidToken).Error
	}

	license, err := service.GetLicenseByLicense(token)
	if err != nil {
		return responder.CreateError(responder.ErrAuthInvalidToken).Error
	}

	if license.ID == 0 {
		responder.WriteJSONResponse(c, responder.CreateError(responder.ErrAuthInvalidToken))
		return responder.CreateError(responder.ErrAuthInvalidToken).Error
	}

	c.Locals("license", license)
	return c.Next()
}
