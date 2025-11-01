package app

import (
	"assignment/config"
	messageprocessor "assignment/internal/pkg/message_processor"
	"assignment/internal/repository/messages"
	"fmt"
	"github.com/redis/go-redis/v9"
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

func (a *APP) GetMessageProcessor() *messageprocessor.Processor {
	return messageprocessor.NewProcessor(a.Config.MessageProcessor)
}
