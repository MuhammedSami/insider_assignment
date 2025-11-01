package messages

import (
	"assignment/internal/repository/models"
	"fmt"
	"gorm.io/gorm"
)

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
