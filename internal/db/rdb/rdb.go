package rdb

import (
	"context"
	"fmt"
	"sync"

	"github.com/kondohiroki/go-boilerplate/config"
	"github.com/kondohiroki/go-boilerplate/internal/logger"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var rdb *redis.Client
var m sync.Mutex
var prefix string
var queuePrefix string

func InitRedisClient(redisConfig config.Redis) error {
	m.Lock()
	defer m.Unlock()

	addr := fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port)
	rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: redisConfig.Password,
		DB:       redisConfig.Database,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return err
	}

	// Set the prefix string
	// for whoever is using AddPrefix() or GetPrefix()
	prefix = config.GetConfig().App.NameSlug

	return nil
}

func GetRedisClient() *redis.Client {
	if rdb == nil {
		m.Lock()
		defer m.Unlock()

		logger.Log.Info("Initializing redis again")
		err := InitRedisClient(config.GetConfig().Redis)
		if err != nil {
			logger.Log.Error("Failed to initialize redis client", zap.Error(err))
		}
		logger.Log.Info("redis initialized")
	}

	return rdb
}

func AddPrefix(key string) string {
	if prefix == "" {
		m.Lock()
		defer m.Unlock()
		prefix = config.GetConfig().App.NameSlug
	}
	return fmt.Sprintf("%s_%s", prefix, key)
}

func AddQueuePrefix(key string) string {
	if queuePrefix == "" {
		m.Lock()
		defer m.Unlock()
		queuePrefix = config.GetConfig().App.NameSlug + "_queue"
	}
	return fmt.Sprintf("%s_%s", queuePrefix, key)
}

func GetPrefix() string {
	return prefix
}

func GetQueuePrefix() string {
	return queuePrefix
}
