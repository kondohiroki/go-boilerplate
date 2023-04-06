package miscellaneous

import (
	"github.com/gofiber/fiber/v2"
)

type MiscellaneousHTTPHandler struct{}

func NewMiscellaneousHTTPHandler() *MiscellaneousHTTPHandler {
	return &MiscellaneousHTTPHandler{}
}

func (m *MiscellaneousHTTPHandler) NotFound(c *fiber.Ctx) error {
	return fiber.NewError(fiber.StatusNotFound, "This is not the api you are looking for, please try again.")
}
