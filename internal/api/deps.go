package api

import repoModels "paribu_assignment/internal/repository/models"

type repo interface {
	GetPendingMessages() ([]repoModels.Message, error)
}
