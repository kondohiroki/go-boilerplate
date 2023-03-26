package healthz

import "github.com/gofiber/fiber/v2"

type HealthzHandler interface {
	Healthz(c *fiber.Ctx) error
}

type healthzContext struct{}

func NewHealthzContext() HealthzHandler {
	return &healthzContext{}
}

func (h *healthzContext) Healthz(c *fiber.Ctx) error {
	c.Status(200).JSON(map[string]any{
		"code":    0,
		"message": "OK",
	})

	return nil
}
