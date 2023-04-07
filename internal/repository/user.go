package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kondohiroki/go-boilerplate/internal/db/model"
)

type UserRepository interface {
	GetUsers(ctx context.Context) ([]model.User, error)
	GetUsersWithPagination(ctx context.Context, limit int, offset int) ([]model.User, error)
	AddUser(ctx context.Context, user model.User) (id int, err error)
}

type UserRepositoryImpl struct {
	pgxPool *pgxpool.Pool
}

func NewUserRepository(pgxPool *pgxpool.Pool) UserRepository {
	return &UserRepositoryImpl{
		pgxPool: pgxPool,
	}
}

func (u *UserRepositoryImpl) GetUsers(ctx context.Context) ([]model.User, error) {
	var users []model.User
	rows, err := u.pgxPool.Query(ctx, "SELECT id, name, email FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var user model.User
		err = rows.Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (u *UserRepositoryImpl) GetUsersWithPagination(ctx context.Context, limit int, offset int) ([]model.User, error) {
	var users []model.User
	rows, err := u.pgxPool.Query(ctx, "SELECT id, name, email FROM users LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var user model.User
		err = rows.Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

// Add user with transaction and return id
func (u *UserRepositoryImpl) AddUser(ctx context.Context, user model.User) (id int, err error) {
	tx, err := u.pgxPool.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(ctx, "INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id", user.Name, user.Email).Scan(&id)
	if err != nil {
		return 0, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return 0, err
	}

	return id, nil
}
