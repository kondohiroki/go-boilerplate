package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/kondohiroki/go-boilerplate/internal/logger"
	"go.uber.org/zap"
)

func Logger() fiber.Handler {
	return func(c *fiber.Ctx) error {
		startTime := time.Now()

		// Continue processing the request
		err := c.Next()

		// Log request
		logger.Log.Info(c.Path(),
			zap.String("request-id", c.Locals("requestid").(string)),
			zap.Int("status", c.Response().StatusCode()),
			zap.String("method", c.Method()),
			zap.String("url", c.Request().URI().String()),
			zap.String("ip", c.IP()),
			zap.String("user-agent", c.Get("User-Agent")),
			zap.String("latency", time.Since(startTime).String()),
		)

		return err
	}
}
