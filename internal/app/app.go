package app

import (
	"fmt"
	"gorm.io/gorm"
	"log"
	"net/http"
	"paribu_assignment/internal/api"
	"paribu_assignment/internal/config"
)

type APP struct {
	API    *api.Manager
	DB     *gorm.DB
	Config *config.Config
}

func NewApp(db *gorm.DB, cfg *config.Config) *APP {
	app := &APP{
		DB:     db,
		Config: cfg,
	}

	app.API = api.NewManager(
		app.GetRedisClient(cfg.Redis),
		app.GetMessagesRepo(),
	)

	return app
}

func (a *APP) Expose() {
	r := a.API.RegisterHandlers()

	err := http.ListenAndServe(fmt.Sprintf(":%d", a.Config.Api.Port), r)
	if err != nil {
		log.Fatalf("failed to serve http err: %+v", err)
	}
}
