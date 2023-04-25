package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/kondohiroki/go-boilerplate/internal/db/rdb"
)

// Set sets a key-value pair with an expiration time.
func Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	key = rdb.AddPrefix(key)
	err := rdb.GetRedisClient().Set(ctx, key, value, expiration).Err()
	if err != nil {
		return fmt.Errorf("failed to set key %s: %w", key, err)
	}
	return nil
}

// Get retrieves the value of a key from Redis.
func Get(ctx context.Context, key string) (string, error) {
	key = rdb.AddPrefix(key)
	val, err := rdb.GetRedisClient().Get(ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("failed to get key %s: %w", key, err)
	}
	return val, nil
}

// Pull retrieves the value of a key from Redis and then deletes the key-value pair.
func Pull(ctx context.Context, key string) (string, error) {
	key = rdb.AddPrefix(key)
	val, err := rdb.GetRedisClient().Get(ctx, key).Result()
	if err != nil {
		return "", fmt.Errorf("failed to get key %s: %w", key, err)
	}

	_, delErr := rdb.GetRedisClient().Del(ctx, key).Result()
	if delErr != nil {
		return "", fmt.Errorf("failed to delete key %s: %w", key, delErr)
	}

	return val, nil
}

// Forever sets the value of a key without an expiration time.
func SetForever(ctx context.Context, key string, value interface{}) error {
	key = rdb.AddPrefix(key)
	err := rdb.GetRedisClient().Set(ctx, key, value, 0).Err()
	if err != nil {
		return fmt.Errorf("failed to set key %s forever: %w", key, err)
	}
	return nil
}

// Delete the key-value pair from Redis.
func Remove(ctx context.Context, key string) error {
	key = rdb.AddPrefix(key)
	_, err := rdb.GetRedisClient().Del(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("failed to forget key %s: %w", key, err)
	}
	return nil
}

// Remove all keys from the current database.
func Flush(ctx context.Context) error {
	key := rdb.AddPrefix("*")
	_, err := rdb.GetRedisClient().Del(ctx, key).Result()
	if err != nil {
		return err
	}
	return nil
}

// Increment increases the integer value of a key by the given increment.
// If the key does not exist, it is set to 0 before performing the operation.
func Increment(ctx context.Context, key string, increment int64) (int64, error) {
	key = rdb.AddPrefix(key)
	val, err := rdb.GetRedisClient().IncrBy(ctx, key, increment).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to increment key %s by %d: %w", key, increment, err)
	}
	return val, nil
}

// Decrement decreases the integer value of a key by the given decrement.
// If the key does not exist, it is set to 0 before performing the operation.
func Decrement(ctx context.Context, key string, decrement int64) (int64, error) {
	key = rdb.AddPrefix(key)
	val, err := rdb.GetRedisClient().DecrBy(ctx, key, decrement).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to decrement key %s by %d: %w", key, decrement, err)
	}
	return val, nil
}

func Remember(ctx context.Context, key string, duration time.Duration, fetchFunc func() ([]byte, error)) ([]byte, error) {
	value, err := Get(ctx, key)
	if err == nil {
		return []byte(value), nil
	}

	data, err := fetchFunc()
	if err != nil {
		return nil, err
	}

	err = Set(ctx, key, data, duration)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func RememberForever(ctx context.Context, key string, fetchFunc func() ([]byte, error)) ([]byte, error) {
	return Remember(ctx, key, 0, fetchFunc)
}
