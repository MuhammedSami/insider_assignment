package tests

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"paribu_assignment/config"
	"paribu_assignment/internal/api/models"
	"paribu_assignment/internal/app"
	"paribu_assignment/internal/storage"
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

	weather := &models.MessageResponse{}
	if err := json.NewDecoder(resp.Body).Decode(weather); err != nil {
		t.Error(err)
	}

	fmt.Println(weather)
}
