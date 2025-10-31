package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
	"net/http"
)

type Manager struct {
	redis       *redis.Client
	messageRepo repo
}

func NewManager(
	redis *redis.Client,
	messageRepo repo,
) *Manager {
	return &Manager{
		redis:       redis,
		messageRepo: messageRepo,
	}
}

func (a *Manager) RegisterHandlers() http.Handler {
	r := chi.NewRouter()

	r.Get("/messages", a.GetMessages)

	return r
}
