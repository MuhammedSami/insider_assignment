package tests

import (
	"assignment/internal/app"
	"assignment/internal/repository/models"
	"assignment/internal/storage"
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestProcessor(t *testing.T) {
	ctx := context.Background()
	cfg := defaultConfig()

	cfg.Message.SendInterval = 2 * time.Second
	cfg.Message.BatchProcessCount = 2

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(`{"sent": true}`))
	}))
	defer server.Close()

	cfg.MessageProcessor.Host = server.URL
	cfg.MessageProcessor.Token = "token"

	db := storage.NewDb(cfg.DB)
	db.Exec("DELETE FROM messages")

	db.Create(&models.Message{
		Content:              "messageContent1",
		RecipientPhoneNumber: "12873",
		Status:               models.StatusPending,
	})
	db.Create(&models.Message{
		Content:              "messageContent2",
		RecipientPhoneNumber: "12874",
		Status:               models.StatusPending,
	})

	app := app.NewApp(ctx, db, cfg)
	assert.True(t, app.API.AutoProcessorRunning, "processor should be running")

	time.Sleep(5 * time.Second)

	var messages []models.Message
	err := db.Find(&messages).Error
	require.NoError(t, err)

	for _, m := range messages {
		assert.Equal(t, models.StatusSent, m.Status, "expected message %s to be sent", m.UUID)
	}
}
