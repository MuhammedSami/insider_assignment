package models

import (
	"time"

	"github.com/google/uuid"
)

type MessageStatus string

const (
	StatusPending       MessageStatus = "pending"
	StatusSent          MessageStatus = "sent"
	StatusFailed        MessageStatus = "failed"
	StatusPermanentFail MessageStatus = "permanent_fail"
)

type Message struct {
	UUID                 uuid.UUID     `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Content              string        `gorm:"type:varchar(1000)"`
	RecipientPhoneNumber string        `gorm:"type:varchar(20)"`
	Status               MessageStatus `gorm:"type:message_status;default:'pending';not null"`
	CreatedAt            time.Time     `gorm:"autoCreateTime"`
}
