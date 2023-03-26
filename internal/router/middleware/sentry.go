package middleware

import (
	"github.com/gofiber/contrib/fibersentry"
	"github.com/gofiber/fiber/v2"
)

func EnhanceSentryEvent() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if hub := fibersentry.GetHubFromContext(c); hub != nil {
			// Set some tags before sending the event to Sentry
			hub.Scope().SetTag("request-id", c.Locals("requestid").(string))
		}
		return c.Next()
	}
}
