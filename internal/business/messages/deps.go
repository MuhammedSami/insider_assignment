package messages

import (
	repoModels "assignment/internal/repository/models"
)

type repo interface {
	GetMessagesByStatuses(limit int, statuses []repoModels.MessageStatus) ([]repoModels.Message, error)
	UpdateStatus(uuid string, status repoModels.MessageStatus) bool
	MessageToRetry(uuid string, maxFailCount int) bool
}
