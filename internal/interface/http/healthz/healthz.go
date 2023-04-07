package healthz

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kondohiroki/go-boilerplate/internal/interface/response"
)

type HealthzHTTPHandler struct{}

func NewHealthzHTTPHandler() *HealthzHTTPHandler {
	return &HealthzHTTPHandler{}
}

func (h *HealthzHTTPHandler) Healthz(c *fiber.Ctx) error {
	c.Status(200).JSON(response.CommonResponse{
		Code:    0,
		Message: "OK",
	})

	return nil
}
