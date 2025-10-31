package api

import (
	"database/sql"
	"github.com/redis/go-redis/v9"
	"net/http"
)

type Handlers struct {
	webhookClient WebhookClient
	db            *sql.DB
	redis         *redis.Client
}

func NewHandlers(
	db *sql.DB,
	redis *redis.Client,
) *Handlers {
	return &Handlers{
		db:    db,
		redis: redis,
	}
}

func (a *Handlers) RegisterHandlers() {
	http.HandleFunc("/weather", a.WeatherHandler)
}
