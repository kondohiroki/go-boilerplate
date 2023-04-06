package router

import (
	"github.com/bytedance/sonic"
	"github.com/gofiber/contrib/fibersentry"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/idempotency"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	httpError "github.com/kondohiroki/go-boilerplate/internal/interface/http/error"
	"github.com/kondohiroki/go-boilerplate/internal/router/middleware"
)

func NewFiberRouter() *fiber.App {
	r := fiber.New(fiber.Config{
		JSONEncoder:           sonic.Marshal,
		JSONDecoder:           sonic.Unmarshal,
		DisableStartupMessage: true,
		EnablePrintRoutes:     false,
		ErrorHandler:          httpError.ErrorHandler,
	})

	// Set up global middleware
	r.Use(cors.New())
	r.Use(requestid.New())
	r.Use(recover.New())
	r.Use(idempotency.New())
	// r.Use(cache.New())
	r.Use(middleware.Logger())
	r.Use(fibersentry.New(fibersentry.Config{
		Repanic:         true,
		WaitForDelivery: true,
	}))
	r.Use(middleware.EnhanceSentryEvent())

	return r
}
