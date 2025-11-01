package tests

import (
	"assignment/config"
	"assignment/internal/api/models"
	"assignment/internal/app"
	"assignment/internal/storage"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetMessages(t *testing.T) {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	cfg.DB.Password = "secret"

	db := storage.NewDb(cfg.DB)

	a := app.NewApp(db, cfg)

	ts := httptest.NewServer(http.HandlerFunc(a.API.GetMessages))

	resp, err := http.Get(fmt.Sprintf("%s/messages", ts.URL))
	if err != nil {
		t.Error(err)
	}

	message := &models.MessageResponse{}
	if err := json.NewDecoder(resp.Body).Decode(message); err != nil {
		t.Error(err)
	}

	fmt.Println(message)
}
