package main

import (
	"assignment/config"
	"assignment/internal/app"
	"assignment/internal/storage"
	"context"
	"log"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	if err := cfg.Validate(); err != nil {
		log.Fatalf("failed to validate config %-+v", err)
	}

	ctx := context.Background()

	db := storage.NewDb(cfg.DB)

	a := app.NewApp(ctx, db, cfg)

	if err := a.ExposeWithGracefulShutDown(ctx); err != nil {
		log.Fatalf("server exited with error: %v", err)
	}
}
