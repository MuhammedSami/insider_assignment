package messageprocessor

import (
	"assignment/config"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const contentType = "application/json"

type MessagePayload struct {
	To      string `json:"to"`
	Content string `json:"content"`
}
type Processor struct {
	cfg config.MessageProcessorAPI
}

func NewProcessor(cfg config.MessageProcessorAPI) *Processor {
	return &Processor{
		cfg: cfg,
	}
}

func (p *Processor) Send(payload MessagePayload) (bool, error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return false, fmt.Errorf("failed to marshal payload: %w", err)
	}

	resp, err := http.Post(fmt.Sprintf("%s/%s", p.cfg.Host, p.cfg.Token), contentType, bytes.NewReader(data))
	if err != nil {
		return false, fmt.Errorf("failed to send message: %w", err)
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode != http.StatusAccepted {
		return false, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return true, nil
}
