package error

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kondohiroki/go-boilerplate/internal/interface/response"
)

// Centralized error handler for all routes
func ErrorHandler(c *fiber.Ctx, err error) error {
	// Status code defaults to 500
	code := fiber.StatusInternalServerError

	// Retrieve request id
	requestID := c.Locals("requestid").(string)

	// Retrieve the custom status code if it's an fiber.*Error
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	// Retrieve the errors from the context
	errors := c.Locals("errors")

	// Handle 400 error
	// if code == fiber.StatusBadRequest {}

	// Handle 401 error
	// if code == fiber.StatusUnauthorized {}

	// Handle 403 error
	// if code == fiber.StatusForbidden {}

	// Handle 404 error
	if code == fiber.StatusNotFound {
		return c.Status(fiber.StatusNotFound).JSON(
			&response.CommonResponse{
				Code:      fiber.StatusNotFound,
				Message:   err.Error(),
				RequestID: requestID,
			},
		)
	}

	// Handle 422 error
	if code == fiber.StatusUnprocessableEntity {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(
			&response.CommonResponse{
				Code:      fiber.StatusUnprocessableEntity,
				Message:   err.Error(),
				Errors:    errors,
				RequestID: requestID,
			},
		)
	}

	// Handle 500 error
	return c.Status(fiber.StatusInternalServerError).JSON(
		&response.CommonResponse{
			Code:      code,
			Message:   err.Error(),
			RequestID: requestID,
		},
	)
}
