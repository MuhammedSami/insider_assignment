package tests

import (
	"assignment/internal/app"
	"assignment/internal/storage"
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProcessor(t *testing.T) {
	cfg := defaultConfig()

	db := storage.NewDb(cfg.DB)

	a := app.NewApp(context.Background(), db, cfg)

	assert.True(t, a.API.AutoProcessorRunning)
}
