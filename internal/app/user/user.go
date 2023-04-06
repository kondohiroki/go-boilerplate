package user

import (
	"github.com/google/uuid"
)

type UserService interface {
	GetUser(input GetUserDTI) (GetUserDTO, error)
	CreateUser(input CreateUserDTI) (CreateUserDTO, error)
}

type userService struct {
	// Add your dependencies here, e.g. database, cache, etc.
}

func NewUserService() UserService {
	return &userService{}
}

type GetUserDTI struct {
	ID int
}

type GetUserDTO struct {
	ID    int
	Name  string
	Email string
}

func (s *userService) GetUser(input GetUserDTI) (GetUserDTO, error) {
	// Replace with actual logic to retrieve the user from the database.
	return GetUserDTO{
		ID:    input.ID,
		Name:  "John Doe",
		Email: "john.doe@example.com",
	}, nil
}

type CreateUserDTI struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

type CreateUserDTO struct {
	ID uuid.UUID `json:"id"`
}

func (s *userService) CreateUser(input CreateUserDTI) (CreateUserDTO, error) {
	// Replace with actual logic to create the user in the database.
	// return random ID
	id := uuid.New()

	return CreateUserDTO{
		ID: id,
	}, nil
}
