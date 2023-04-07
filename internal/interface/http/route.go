package httpapi

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kondohiroki/go-boilerplate/internal/app/user"
	httpHealthz "github.com/kondohiroki/go-boilerplate/internal/interface/http/healthz"
	httpMiscellaneous "github.com/kondohiroki/go-boilerplate/internal/interface/http/miscellaneous"
	httpUser "github.com/kondohiroki/go-boilerplate/internal/interface/http/user"
)

// ====================================================
// =================== DEFINE ROUTE ===================
// ====================================================
func RegisterRoute(r *fiber.App) {
	api := r.Group("/api")
	v1 := api.Group("/v1")

	// Healthz API
	healthzAPI := v1.Group("/healthz")
	healthzHandler := httpHealthz.NewHealthzHTTPHandler()
	healthzAPI.Get("/", healthzHandler.Healthz)

	// User API
	userAPI := v1.Group("/users")
	userService := user.NewUserService()
	userHandler := httpUser.NewUserHTTPHandler(userService)
	userAPI.Get("/", userHandler.GetUsers)
	userAPI.Get("/:id", userHandler.GetUserByID)
	userAPI.Post("/", userHandler.CreateUser)

	// Error Case Handler
	miscellaneousHandler := httpMiscellaneous.NewMiscellaneousHTTPHandler()
	r.All("*", miscellaneousHandler.NotFound)
}
