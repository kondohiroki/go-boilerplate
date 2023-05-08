package user

import (
	"context"

	"github.com/kondohiroki/go-boilerplate/internal/db/model"
	"github.com/kondohiroki/go-boilerplate/internal/repository"
	"github.com/kondohiroki/go-boilerplate/pkg/exception"
)

type UserApp interface {
	GetUsers(ctx context.Context) ([]GetUserDTO, error)
	GetUserByID(ctx context.Context, input GetUserDTI) (GetUserDTO, error)
	CreateUser(ctx context.Context, input CreateUserDTI) (CreateUserDTO, error)
}

type userApp struct {
	Repo *repository.Repository
}

func NewUserApp(repo *repository.Repository) UserApp {
	return &userApp{
		Repo: repo,
	}
}

type GetUserDTI struct {
	ID string
}

type GetUserDTO struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (s *userApp) GetUsers(ctx context.Context) ([]GetUserDTO, error) {
	users, err := s.Repo.User.GetUsers(ctx)
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

func (s *userApp) GetUserByID(ctx context.Context, input GetUserDTI) (GetUserDTO, error) {
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

func (s *userApp) CreateUser(ctx context.Context, input CreateUserDTI) (CreateUserDTO, error) {
	// Ensure email is not already taken
	isUserEmailExist, err := s.Repo.User.IsUserEmailExist(ctx, input.Email)
	if err != nil {
		return CreateUserDTO{}, err
	}
	if isUserEmailExist {
		return CreateUserDTO{}, exception.UserEmailAlreadyTakenError
	}

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
