package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID         int64          `json:"-" gorm:"primaryKey;autoIncrement"`
	PublicID   uuid.UUID      `json:"id" gorm:"type:uuid;uniqueIndex;not null;default:gen_random_uuid()"`
	Name       string         `json:"name"  gorm:"size:100;not null"`
	Email      string         `json:"email"  gorm:"size:100;unique;not null"`
	Password   string         `json:"password,omitempty"  gorm:"size:255;not null"`
	CreatedAt  time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
	UserMatrix *UserMatrix    `gorm:"foreignKey:UserID"`
}

type UserCreateRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email" `
	Password string `json:"password,omitempty" validate:"required,min=6"`
}

type UserMe struct {
	PublicID string `json:"pid" `
	Name     string `json:"username"`
	Email    string `json:"email"`
}

type UserMeResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    UserMe `json:"data"`
}
