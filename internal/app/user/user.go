package user

import (
	"context"

	"github.com/kondohiroki/go-boilerplate/internal/db/model"
	"github.com/kondohiroki/go-boilerplate/internal/repository"
)

type UserService interface {
	GetUsers(ctx context.Context) ([]GetUserDTO, error)
	GetUserByID(ctx context.Context, input GetUserDTI) (GetUserDTO, error)
	CreateUser(ctx context.Context, input CreateUserDTI) (CreateUserDTO, error)
}

type userService struct {
	Repo repository.Repository
}

func NewUserService() UserService {
	return &userService{
		Repo: *repository.NewRepository(),
	}
}

type GetUserDTI struct {
	ID int
}

type GetUserDTO struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (s *userService) GetUsers(ctx context.Context) ([]GetUserDTO, error) {
	users, err := s.Repo.User.GetUsersWithPagination(ctx, 10, 0)
	if err != nil {
		return nil, err
	}

	var usersDTO []GetUserDTO
	for _, user := range users {
		usersDTO = append(usersDTO, GetUserDTO{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		})
	}

	return usersDTO, nil
}

func (s *userService) GetUserByID(ctx context.Context, input GetUserDTI) (GetUserDTO, error) {
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
	ID int `json:"id"`
}

func (s *userService) CreateUser(ctx context.Context, input CreateUserDTI) (CreateUserDTO, error) {

	id, err := s.Repo.User.AddUser(ctx, model.User{
		Name:  input.Name,
		Email: input.Email,
	})

	if err != nil {
		return CreateUserDTO{}, err
	}

	return CreateUserDTO{
		ID: id,
	}, nil
}
