package messages

import (
	repoModels "assignment/internal/repository/models"
	"context"
)

type repo interface {
	MessageToRetry(uuid string, maxFailCount int) bool
	UpdateStatus(uuid string, status repoModels.MessageStatus) bool
	GetMessagesByStatusesWithLock(ctx context.Context, limit int) ([]repoModels.Message, error)
	GetMessagesByStatuses(limit int, statuses []repoModels.MessageStatus) ([]repoModels.Message, error)
}
