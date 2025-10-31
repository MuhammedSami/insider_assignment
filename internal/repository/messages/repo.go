package messages

import (
	"fmt"
	"gorm.io/gorm"
	"paribu_assignment/internal/repository/models"
)

type MessageRepo struct {
	db *gorm.DB
}

func NewMessageRepo(db *gorm.DB) *MessageRepo {
	return &MessageRepo{
		db: db,
	}
}

func (m *MessageRepo) GetPendingMessages() ([]models.Message, error) {
	var messages []models.Message

	// we dont have any limit param but we could have :)
	result := m.db.Where("status = ?", models.StatusPending).Limit(20).Find(&messages)
	if result.Error != nil {
		return nil, fmt.Errorf("failed to retrieve pending messages %w", result.Error)
	}

	return messages, nil
}
