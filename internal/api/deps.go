package api

import (
	repoModels "assignment/internal/repository/models"
	"context"
)

type repo interface {
	GetMessagesByStatuses(limit int, statuses []repoModels.MessageStatus) ([]repoModels.Message, error)
}

type autoMessageProcessor interface {
	Process(ctx context.Context) error
}
