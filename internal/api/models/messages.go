package models

import repoModels "paribu_assignment/internal/repository/models"

type MessageResponse struct {
	Count    int                  `json:"count"`
	Messages []repoModels.Message `json:"messages"`
}
