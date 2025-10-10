package model

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID             int64        `json:"-" gorm:"primaryKey;autoIncrement"`
	PublicID       uuid.UUID    `json:"id" gorm:"type:uuid;uniqueIndex;not null;default:gen_random_uuid()"`
	ConversationID int64        `json:"-" gorm:"index;not null"`
	Conversation   Conversation `json:"-" gorm:"constraint:OnDelete:CASCADE"`
	Role           string       `json:"role" gorm:"type:varchar(16);not null"`
	Content        string       `json:"content" gorm:"type:text;not null"`
	RefID          *int64       `json:"refId,omitempty" `
	CreatedAt      time.Time    `json:"createdAt" gorm:"not null;default:now()"`
	UpdatedAt      time.Time    `json:"updatedAt" gorm:"not null;autoUpdateTime"`
}

type MessageListRes struct {
	Code            int        `json:"code"`
	Message         string     `json:"message"`
	ConversationPID string     `json:"conversationPid"`
	Data            *[]Message `json:"data"`
}

const (
	ROLE_ASSISTANT = "assistant"
	ROLE_SYSTEM    = "system"
	ROLE_USER      = "user"
)
