package api

import (
	"encoding/json"
	"net/http"
	"paribu_assignment/internal/api/models"
)

func (a *Manager) GetMessages(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	messages, err := a.messageRepo.GetPendingMessages()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := models.MessageResponse{
		Count:    len(messages),
		Messages: messages,
	}

	json.NewEncoder(w).Encode(resp)
}
