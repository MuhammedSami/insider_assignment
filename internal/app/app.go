package app

import (
	"assignment/config"
	"assignment/internal/api"
	"assignment/internal/business/messages"
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type APP struct {
	API    api.Manager
	DB     *gorm.DB
	Config *config.Config
}

func NewApp(ctx context.Context, db *gorm.DB, cfg *config.Config) *APP {
	app := &APP{
		DB:     db,
		Config: cfg,
	}

	messageRepo := app.GetMessagesRepo()
	redisClient := app.GetRedisClient(cfg.Redis)

	ctxWithCancel, cancel := context.WithCancel(ctx)

	autoMessageProcessor := messages.NewAuthMessageProcessor(
		cfg,
		messageRepo,
		app.GetMessageProcessor(),
		redisClient,
	)

	app.API = *api.NewManager(
		redisClient,
		messageRepo,
		autoMessageProcessor,
	)

	err := autoMessageProcessor.Process(ctxWithCancel)
	if err != nil {
		log.Fatalf("failed to start auto process err. %+v", err)
	}

	app.API.AutoProcessorCancelFn = cancel
	app.API.AutoProcessorRunning = true

	return app
}

func (a *APP) ExposeWithGracefulShutDown(ctx context.Context) error {
	r := a.API.RegisterHandlers()

	// we could read from env that is right :)
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", a.Config.Api.Port),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	serverError := make(chan error, 1)

	go func() {
		log.Printf("server is running on http://localhost%s", server.Addr)
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			serverError <- err
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverError:
		log.Printf("server error: %v", err)
	case sig := <-stop:
		log.Printf("received shutdown signal: %v", sig)
	}

	log.Println("server is shutting down...")

	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown error: %+v", err)
	}

	// shutdown DB as well, no need for hanging connection on network layer
	sqlDB, err := a.DB.DB()
	if err != nil {
		log.Printf("failed to get sql.DB: %v", err)
	} else {
		if err := sqlDB.Close(); err != nil {
			log.Printf("failed to close db: %v", err)
		}
	}

	log.Println("server exited properly")

	return nil
}
