package tests

import (
	"assignment/config"
	"assignment/internal/app"
	"assignment/internal/storage"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestProcessor(t *testing.T) {
	cfg, err := config.NewConfig()
	require.NoError(t, err)

	cfg.DB.Password = "secret"

	db := storage.NewDb(cfg.DB)

	a := app.NewApp(context.Background(), db, cfg)

	assert.True(t, a.API.AutoProcessorRunning)
}
