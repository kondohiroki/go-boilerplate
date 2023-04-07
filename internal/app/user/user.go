package user

import (
	"github.com/google/uuid"
)

type UserService interface {
	GetUsers() ([]GetUserDTO, error)
	GetUserByID(input GetUserDTI) (GetUserDTO, error)
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
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (s *userService) GetUsers() ([]GetUserDTO, error) {
	// Replace with actual logic to retrieve the users from the database.
	return []GetUserDTO{
		{
			ID:    1,
			Name:  "John Doe",
			Email: "john@gmail.com",
		},
		{
			ID:    2,
			Name:  "Lucy",
			Email: "lucy@gmail.com",
		},
	}, nil
}

func (s *userService) GetUserByID(input GetUserDTI) (GetUserDTO, error) {
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
