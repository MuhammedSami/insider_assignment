//nolint:errcheck
package api

import (
	apiModels "assignment/internal/api/models"
	repoModels "assignment/internal/repository/models"
	"context"
	"encoding/json"
	"net/http"
)

func (a *Manager) GetMessages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	defaultLimit := 20 // we could receive it as query param but let's skip it

	messages, err := a.messageRepo.GetMessagesByStatuses(defaultLimit, []repoModels.MessageStatus{repoModels.StatusSent}) // we could use a general entity to not import repo but I am skipping for now
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var apiMessages []apiModels.Message

	for _, m := range messages {
		// we dont return failedCount bec it is not relevant
		apiMessages = append(apiMessages, apiModels.Message{
			Id:                   m.UUID.String(),
			Content:              m.Content,
			RecipientPhoneNumber: m.RecipientPhoneNumber,
			Status:               m.Status.ToString(),
			CreatedAt:            m.CreatedAt.String(),
		})
	}

	resp := apiModels.MessageResponse{
		Count:    len(messages),
		Messages: apiMessages,
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message":"failed to encode response"}`))

		return
	}
}

func (a *Manager) StartProcessor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if a.AutoProcessorRunning {
		w.WriteHeader(http.StatusConflict)
		w.Write([]byte(`{"message":processor already running"}`))
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	a.AutoProcessorCancelFn = cancel
	a.AutoProcessorRunning = true

	err := a.AutoMessageProcessor.Process(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message":"failed to start processor"}`))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"processor started"}`))
}

func (a *Manager) StopProcessor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if !a.AutoProcessorRunning {
		w.WriteHeader(http.StatusConflict) // we could return 200 as well it makes sense but I think it depends on who is using this API
		w.Write([]byte(`{"message":"processor not running"}`))
		return
	}

	a.AutoProcessorCancelFn()
	a.AutoProcessorRunning = false

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message":"processor stopped"}`))
}
