package app

import (
	"assignment/config"
	"assignment/internal/api"
	"assignment/internal/business/messages"
	"context"
	"fmt"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type APP struct {
	API    api.Manager
	DB     *gorm.DB
	Config *config.Config
}

func NewApp(db *gorm.DB, cfg *config.Config) *APP {
	app := &APP{
		DB:     db,
		Config: cfg,
	}

	messageRepo := app.GetMessagesRepo()
	redisClient := app.GetRedisClient(cfg.Redis)

	app.API = *api.NewManager(
		redisClient,
		messageRepo,
	)

	ctx, cancel := context.WithCancel(context.Background())

	autoMessageProcessor := messages.NewAuthMessageProcessor(
		cfg,
		messageRepo,
		app.GetMessageProcessor(),
		redisClient,
	)

	err := autoMessageProcessor.Process(ctx)
	if err != nil {
		log.Fatalf("failed to start auto process err. %+v", err)
	}

	app.API.AutoProcessorCancelFn = cancel
	app.API.AutoProcessorRunning = true

	return app
}

func (a *APP) Expose() {
	r := a.API.RegisterHandlers()

	err := http.ListenAndServe(fmt.Sprintf(":%d", a.Config.Api.Port), r)
	if err != nil {
		log.Fatalf("failed to serve http err: %+v", err)
	}
}
