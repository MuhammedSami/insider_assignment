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
	AutoMessageProcessor  autoMessageProcessor
	AutoProcessorCancelFn context.CancelFunc
	AutoProcessorRunning  bool
}

func NewManager(
	redis *redis.Client,
	messageRepo repo,
	autoMessageProcessor autoMessageProcessor,
) *Manager {
	return &Manager{
		redis:                redis,
		messageRepo:          messageRepo,
		AutoMessageProcessor: autoMessageProcessor,
	}
}

func (a *Manager) RegisterHandlers() http.Handler {
	r := chi.NewRouter()

	r.Get("/messages/sent", a.GetMessages)
	r.Post("/processor/start", a.StartProcessor)
	r.Post("/processor/stop", a.StopProcessor)

	return r
}
