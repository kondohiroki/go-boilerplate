package httpapi

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kondohiroki/go-boilerplate/internal/app/healthz"
	"github.com/kondohiroki/go-boilerplate/internal/app/miscellaneous"
)

// ====================================================
// =================== DEFINE ROUTE ===================
// ====================================================
func RegisterRoute(r *fiber.App) {
	api := r.Group("/api")
	// v1 := api.Group("/v1")

	// Healthz Handler
	healthzHandler := healthz.NewHealthzContext()
	api.Get("/healthz", healthzHandler.Healthz)

	// Error Case Handler
	MiscellaneousHandler := miscellaneous.NewMiscellaneousHandler()
	r.All("*", MiscellaneousHandler.NotFound)
}
