package api

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
	"net/http"
)

type Manager struct {
	redis                 *redis.Client
	messageRepo           repo
	AutoProcessorCancelFn context.CancelFunc
	AutoProcessorRunning  bool
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

	r.Get("/messages/sent", a.GetMessages)
	r.Post("/processor/start", a.StartStopProcessor)
	r.Post("/processor/stop", a.StartStopProcessor)

	return r
}
