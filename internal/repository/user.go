package repository

import (
	"context"
	"time"

	"github.com/bytedance/sonic"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kondohiroki/go-boilerplate/internal/db/model"
	"github.com/kondohiroki/go-boilerplate/internal/db/rdb"
	"github.com/redis/go-redis/v9"
)

type UserRepository interface {
	GetUsers(ctx context.Context) ([]model.User, error)
	GetUsersWithPagination(ctx context.Context, limit int, offset int) ([]model.User, error)
	AddUser(ctx context.Context, user model.User) (id int, err error)
}

type UserRepositoryImpl struct {
	pgxPool     *pgxpool.Pool
	redisClient *redis.Client
}

func NewUserRepository(pgxPool *pgxpool.Pool, redisClient *redis.Client) UserRepository {
	return &UserRepositoryImpl{
		pgxPool:     pgxPool,
		redisClient: redisClient,
	}
}

func (u *UserRepositoryImpl) GetUsers(ctx context.Context) ([]model.User, error) {
	key := "users"

	data, err := rdb.Remember(ctx, key, 10*time.Minute, func() ([]byte, error) {
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

		// Serialize users to bytes using Sonic
		userBytes, err := sonic.Marshal(users)
		if err != nil {
			return nil, err
		}

		return userBytes, nil
	})

	if err != nil {
		return nil, err
	}

	// Deserialize data to []model.User
	var users []model.User
	err = sonic.Unmarshal(data, &users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

// Work in progress
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

	// Delete cache
	err = rdb.Remove(ctx, "users")

	return id, nil
}
