package api

import (
	apiModels "assignment/internal/api/models"
	repoModels "assignment/internal/repository/models"
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

	json.NewEncoder(w).Encode(resp)
}

func (a *Manager) StartStopProcessor(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if a.AutoProcessorRunning {
		a.AutoProcessorCancelFn()
	}
}
