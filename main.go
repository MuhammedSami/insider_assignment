package main

import (
	"log"
	"net/http"
	"paribu_assignment/internal/api"
	"paribu_assignment/internal/config"
	"paribu_assignment/internal/storage"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := storage.NewDb(cfg.DB)
	if err != nil {
		log.Fatalf("failed to connect DB: %v", err)
	}
	defer db.Close()

	redisClient := storage.GetRedisClient(cfg.Redis)

	a := api.NewHandlers(db, redisClient)
	a.RegisterHandlers()

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("failed to serve http err: %+v", err)
	}
}
