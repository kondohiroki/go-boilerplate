package user

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/kondohiroki/go-boilerplate/internal/app/user"
	"github.com/kondohiroki/go-boilerplate/internal/interface/response"
	"github.com/kondohiroki/go-boilerplate/internal/interface/validation"
	"github.com/kondohiroki/go-boilerplate/pkg/exception"
)

type UserHTTPHandler struct {
	app user.UserApp
}

func NewUserHTTPHandler(app user.UserApp) *UserHTTPHandler {
	return &UserHTTPHandler{app: app}
}

// Write me GetUsers function
func (h *UserHTTPHandler) GetUsers(c *fiber.Ctx) error {
	dtos, err := h.app.GetUsers(c.Context())
	if err != nil {
		return err
	}

	return c.JSON(response.CommonResponse{
		ResponseCode:    0,
		ResponseMessage: "OK",
		Data:            dtos,
	})
}

func (h *UserHTTPHandler) GetUserByID(c *fiber.Ctx) error {
	id := c.Params("id")

	dti := user.GetUserDTI{ID: id}
	dto, err := h.app.GetUserByID(c.Context(), dti)
	if err != nil {
		return err
	}

	return c.JSON(response.CommonResponse{
		ResponseCode:    0,
		ResponseMessage: "OK",
		Data:            dto,
	})
}

func (h *UserHTTPHandler) CreateUser(c *fiber.Ctx) error {
	var req user.CreateUserDTI

	// Parse the request body
	if err := c.BodyParser(&req); err != nil {
		return exception.InvalidRequestBodyError
	}

	// Validate the request body
	v, _ := validation.GetValidator()
	if err := v.Struct(req); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			return exception.NewValidationFailedErrors(validationErrs)
		}
	}

	// Process the business logic
	dto, err := h.app.CreateUser(c.Context(), user.CreateUserDTI{
		Name:  req.Name,
		Email: req.Email,
	})

	if err != nil {
		return err
	}

	return c.Status(http.StatusCreated).JSON(response.CommonResponse{
		ResponseCode:    0,
		ResponseMessage: "OK",
		Data:            dto,
	})
}
