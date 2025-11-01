package messageprocessor

import (
	"assignment/config"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const contentType = "application/json"

type MessagePayload struct {
	To      string `json:"to"`
	Content string `json:"content"`
}
type Processor struct {
	cfg    config.MessageProcessorAPI
	client *http.Client
}

func NewProcessor(cfg config.MessageProcessorAPI) *Processor {
	return &Processor{
		cfg: cfg,
		// we could use a wrapper in case where we have multiple calls to different endpoints with different headers, methods etc..
		// a builder pattern, but just putting this here for showcasing the importance of http.Client with correct config.
		// usually this will come from some config for Transport control over http.
		client: &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:          100,
				MaxConnsPerHost:       100,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   5 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
			},
		},
	}
}

func (p *Processor) Send(ctx context.Context, payload MessagePayload) (bool, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return false, fmt.Errorf("failed to marshal payload: %w", err)
	}

	url := fmt.Sprintf("%s/%s", p.cfg.Host, p.cfg.Token)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(data))
	if err != nil {
		return false, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return false, fmt.Errorf("failed to send message: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return true, nil
}
