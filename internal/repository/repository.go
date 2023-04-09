package repository

import (
	"github.com/kondohiroki/go-boilerplate/internal/db/pgx"
	"github.com/kondohiroki/go-boilerplate/internal/db/rdb"
)

type Repository struct {
	User UserRepository
}

func NewRepository() *Repository {
	pgxPool := pgx.GetPgxPool()
	redisClient := rdb.GetRedisClient()

	return &Repository{
		User: NewUserRepository(pgxPool, redisClient),
	}
}
