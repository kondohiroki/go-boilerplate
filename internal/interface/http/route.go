package httpapi

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kondohiroki/go-boilerplate/internal/app/queue"
	"github.com/kondohiroki/go-boilerplate/internal/app/user"

	httpHealthz "github.com/kondohiroki/go-boilerplate/internal/interface/http/healthz"
	httpMiscellaneous "github.com/kondohiroki/go-boilerplate/internal/interface/http/miscellaneous"
	httpQueue "github.com/kondohiroki/go-boilerplate/internal/interface/http/queue"
	httpUser "github.com/kondohiroki/go-boilerplate/internal/interface/http/user"
)

// ====================================================
// =================== DEFINE ROUTE ===================
// ====================================================
func RegisterRoute(r *fiber.App) {
	api := r.Group("/api")
	v1 := api.Group("/v1")

	// Healthz API
	healthzAPI := api.Group("/healthz")
	healthzHandler := httpHealthz.NewHealthzHTTPHandler()
	healthzAPI.Get("/", healthzHandler.Healthz)

	// User API
	userAPI := v1.Group("/users")
	userApp := user.NewUserApp()
	userHandler := httpUser.NewUserHTTPHandler(userApp)
	userAPI.Get("/", userHandler.GetUsers)
	userAPI.Get("/:id", userHandler.GetUserByID)
	userAPI.Post("/", userHandler.CreateUser)

	// Queue API
	queueAPI := v1.Group("/queues")
	queueApp := queue.NewQueueApp()
	queueHandler := httpQueue.NewQueueHTTPHandler(queueApp)
	queueAPI.Get("/", queueHandler.GetQueues)
	// queueAPI.Get("/:key", queueHandler.GetQueueByKey)

	// Error Case Handler
	miscellaneousHandler := httpMiscellaneous.NewMiscellaneousHTTPHandler()
	r.All("*", miscellaneousHandler.NotFound)
}
