package storage

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"paribu_assignment/config"
)

func GetRedisClient(cfg config.RedisConn) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       0,
	})
}
