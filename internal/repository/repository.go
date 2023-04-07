package repository

import "github.com/kondohiroki/go-boilerplate/internal/db"

type Repository struct {
	User UserRepository
}

func NewRepository() *Repository {
	pgxPool := db.GetPgxPool()

	return &Repository{
		User: NewUserRepository(pgxPool),
	}
}
