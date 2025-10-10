package middleware

import (
	"cms-octo-chat-api/model"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type (
	BaseMiddleware struct {
		Env  *model.EnvVar
		Logs *logrus.Logger
		DB   *gorm.DB
	}
)

func NewBaseMiddleware(m BaseMiddleware) *BaseMiddleware {
	return &BaseMiddleware{
		Env:  m.Env,
		Logs: m.Logs,
		DB:   m.DB,
	}
}
