package app

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"paribu_assignment/internal/config"
	"paribu_assignment/internal/repository/messages"
)

func (a *APP) GetRedisClient(cfg config.RedisConn) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       0,
	})
}

func (a *APP) GetMessagesRepo() *messages.MessageRepo {
	return messages.NewMessageRepo(a.DB)
}
