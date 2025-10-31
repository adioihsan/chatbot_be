package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Conversation struct {
	ID        int64          `json:"-" gorm:"primaryKey;autoIncrement"` // hide internal ID from API
	PublicID  uuid.UUID      `json:"pid" gorm:"type:uuid;uniqueIndex;not null;default:gen_random_uuid()"`
	Title     string         `json:"title" gorm:"type:varchar(160);not null;default:'Untitled chat'"`
	UserID    int64          `json:"-" gorm:"index"`
	CreatedAt time.Time      `json:"createdAt" gorm:"not null;default:now()"` // or `autoCreateTime`
	UpdatedAt time.Time      `json:"updatedAt" gorm:"not null;autoUpdateTime"`
	DeletedAt gorm.DeletedAt ` json:"deletedAt" gorm:"index"`
	Messages  []Message      `json:"messages,omitempty" gorm:"constraint:OnDelete:CASCADE"`
}

type ConversationCreateReq struct {
	Title string `json:"title" validate:"required,max=64"`
}

type ConversationCreateRes struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Data    *Conversation `json:"data"`
}

type ConversationListRes struct {
	Code    int            `json:"code"`
	Message string         `json:"message"`
	LastPid *uuid.UUID     `json:"lastPid,omitempty"`
	Data    []Conversation `json:"data"`
}

type ConversationRenameReq struct {
	Title string `json:"title" validate:"required,max=64"`
}
