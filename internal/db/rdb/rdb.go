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

var rdb redis.Cmdable
var m sync.Mutex
var prefix string
var queuePrefix string

type RedisCredentials struct {
	Password string
	Database int
}

func InitRedisClient(redisConfigs []config.Redis) error {
	m.Lock()
	defer m.Unlock()

	// Prepare a list of Redis addresses
	// Prepare a list of Redis addresses and a map of their corresponding credentials
	var addrs []string
	creds := make(map[string]RedisCredentials)
	for _, redisConfig := range redisConfigs {
		addr := fmt.Sprintf("%s:%d", redisConfig.Host, redisConfig.Port)
		addrs = append(addrs, addr)
		creds[addr] = RedisCredentials{
			Password: redisConfig.Password,
			Database: redisConfig.Database,
		}
	}

	if len(addrs) == 1 {
		rdb = redis.NewClient(&redis.Options{
			Addr:     addrs[0],
			Password: creds[addrs[0]].Password,
			DB:       creds[addrs[0]].Database,
		})
	} else {
		rdb = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs: addrs,
			NewClient: func(opt *redis.Options) *redis.Client {
				cred := creds[opt.Addr]
				opt.Password = cred.Password
				opt.DB = cred.Database

				return redis.NewClient(opt)
			},
		})
	}

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return err
	}

	// Set the prefix string
	// for whoever is using AddPrefix() or GetPrefix()
	prefix = config.GetConfig().App.NameSlug
	queuePrefix = config.GetConfig().App.NameSlug + "_queue"

	return nil
}

func GetRedisClient() redis.Cmdable {
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
