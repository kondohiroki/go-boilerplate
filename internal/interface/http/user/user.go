package user

import (
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/kondohiroki/go-boilerplate/internal/app/user"
	"github.com/kondohiroki/go-boilerplate/internal/interface/response"
	"github.com/kondohiroki/go-boilerplate/internal/interface/validation"
)

type UserHTTPHandler struct {
	service user.UserService
}

func NewUserHTTPHandler(service user.UserService) *UserHTTPHandler {
	return &UserHTTPHandler{service: service}
}

func (h *UserHTTPHandler) GetUsers(c *fiber.Ctx) error {
	dtos, err := h.service.GetUsers(c.Context())
	if err != nil {
		return err
	}

	return c.JSON(response.CommonResponse{
		Code:    0,
		Message: "OK",
		Data:    dtos,
	})
}

func (h *UserHTTPHandler) GetUserByID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "Invalid ID")
	}

	dti := user.GetUserDTI{ID: id}
	dto, err := h.service.GetUserByID(c.Context(), dti)
	if err != nil {
		return err
	}

	return c.JSON(response.CommonResponse{
		Code:    0,
		Message: "OK",
		Data:    dto,
	})
}

func (h *UserHTTPHandler) CreateUser(c *fiber.Ctx) error {
	var req user.CreateUserDTI

	// Parse the request body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).SendString("Invalid request body")
	}

	// Validate the request body
	v, _ := validation.GetValidator()
	if err := v.Struct(req); err != nil {
		errors := validation.GetValidationErrors(err.(validator.ValidationErrors))
		c.Context().SetUserValue("errors", errors)
		return fiber.NewError(fiber.StatusUnprocessableEntity, "Request body is not valid format or missing required fields")
	}

	// Process the business logic
	dto, err := h.service.CreateUser(c.Context(), user.CreateUserDTI{
		Name:  req.Name,
		Email: req.Email,
	})

	if err != nil {
		return err
	}

	return c.JSON(response.CommonResponse{
		Code:    0,
		Message: "OK",
		Data:    dto,
	})
}
