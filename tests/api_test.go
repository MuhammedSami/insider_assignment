package tests

import (
	"assignment/config"
	"assignment/internal/app"
	"assignment/internal/storage"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPIEndpoints(t *testing.T) {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	cfg.DB.Password = "secret"
	db := storage.NewDb(cfg.DB)
	a := app.NewApp(db, cfg)

	ts := httptest.NewServer(a.API.RegisterHandlers())
	defer ts.Close()

	tests := []struct {
		name           string
		method         string
		endpoint       string
		expectedStatus int
	}{
		{
			name:           "Get sent messages",
			method:         http.MethodGet,
			endpoint:       "/messages/sent",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Stop processor",
			method:         http.MethodPost,
			endpoint:       "/processor/stop",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Start processor",
			method:         http.MethodPost,
			endpoint:       "/processor/start",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Start processor when already running",
			method:         http.MethodPost,
			endpoint:       "/processor/start",
			expectedStatus: http.StatusConflict,
		},
		{
			name:           "Stop processor when already stopped",
			method:         http.MethodPost,
			endpoint:       "/processor/stop",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Stop processor when already stopped",
			method:         http.MethodPost,
			endpoint:       "/processor/stop",
			expectedStatus: http.StatusConflict,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var resp *http.Response
			var err error

			switch tt.method {
			case http.MethodGet:
				resp, err = http.Get(ts.URL + tt.endpoint)
			case http.MethodPost:
				resp, err = http.Post(ts.URL+tt.endpoint, "application/json", nil)
			}

			if err != nil {
				t.Fatalf("request failed: %v", err)
			}
			defer resp.Body.Close() //nolint:errcheck

			assert.Equal(t, tt.expectedStatus, resp.StatusCode)
		})
	}
}
