package miscellaneous

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kondohiroki/go-boilerplate/internal/interface/httpapi/response"
)

type MiscellaneousHandler interface {
	NotFound(c *fiber.Ctx) error
}

type miscellaneousContext struct{}

func NewMiscellaneousHandler() MiscellaneousHandler {
	return &miscellaneousContext{}
}

func (m *miscellaneousContext) NotFound(c *fiber.Ctx) error {
	c.Status(fiber.StatusNotFound).JSON(response.Response{
		Code:    0,
		Message: "This is not the api you are looking for, please try again.",
		Data:    nil,
	})
	return nil
}
