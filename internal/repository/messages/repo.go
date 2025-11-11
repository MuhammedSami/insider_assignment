package messages

import (
	"assignment/internal/repository/models"
	"context"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const defaultLimit = 20

type MessageRepo struct {
	db *gorm.DB
}

func NewMessageRepo(db *gorm.DB) *MessageRepo {
	return &MessageRepo{
		db: db,
	}
}

func (m *MessageRepo) GetMessagesByStatuses(limit int, statuses []models.MessageStatus) ([]models.Message, error) {
	var messages []models.Message

	query := m.db.Where("status IN ?", statuses)
	if limit > 0 {
		query = query.Limit(limit)
	}

	result := query.Find(&messages)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to retrieve messages with statuses %v: %w", statuses, result.Error)
	}

	return messages, nil
}

func (m *MessageRepo) GetMessagesByStatusesWithLock(ctx context.Context, limit int) ([]models.Message, error) {
	var messages []models.Message

	if limit == 0 {
		limit = defaultLimit
	}

	tx := m.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	query := tx.Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).
		Where("status IN ?", []models.MessageStatus{models.StatusPending, models.StatusFailed}).
		Limit(limit).
		Find(&messages)
	if query.Error != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to fetch messages with lock, err:%+v", query.Error)
	}

	if len(messages) == 0 {
		return nil, nil
	}

	ids := make([]uuid.UUID, len(messages))
	for i, msg := range messages {
		ids[i] = msg.UUID
	}

	if err := tx.Model(models.Message{}).Where("uuid IN ?", ids).Update("status", models.StatusProcessing).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to lock messages err:%+v", err)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction, err:%+v", err)
	}

	return messages, nil
}

func (m *MessageRepo) UpdateStatus(uuid string, status models.MessageStatus) bool {
	result := m.db.Model(&models.Message{}).
		Where("uuid = ?", uuid).
		Update("status", status)
	if result.Error != nil {
		return false
	}

	return result.RowsAffected > 0
}

func (m *MessageRepo) MessageToRetry(uuid string, maxFailCount int) bool {
	result := m.db.Model(&models.Message{}).
		Where("uuid = ?", uuid).
		Updates(map[string]interface{}{
			"failed_count": gorm.Expr("failed_count + 1"),
			"status": gorm.Expr(`
            CASE 
                WHEN failed_count + 1 >= ` + fmt.Sprintf("%d", maxFailCount) + ` THEN 'permanent_fail'::message_status
                ELSE 'failed'::message_status
            END
        `),
		})

	return result.RowsAffected > 0
}
