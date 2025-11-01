package main

import (
	"assignment/config"
	"assignment/internal/app"
	"assignment/internal/storage"
	"log"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	if err := cfg.Validate(); err != nil {
		log.Fatalf("failed to validate config %w", err)
	}

	db := storage.NewDb(cfg.DB)

	a := app.NewApp(db, cfg)
	a.Expose()
}
