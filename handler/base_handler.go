package handler

import (
	"cms-octo-chat-api/model"
	"cms-octo-chat-api/repository"

	openai "github.com/openai/openai-go/v3"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type (
	BaseHandler struct {
		Env    *model.EnvVar
		Logs   *logrus.Logger
		DB     *gorm.DB
		Repo   *repository.BaseRepository
		OpenAI *openai.Client
	}
)

func NewBaseHandler(obj BaseHandler) *BaseHandler {

	r := repository.NewBaseRepository(repository.BaseRepository{
		Env:  obj.Env,
		Logs: obj.Logs,
		DB:   obj.DB,
	})

	return &BaseHandler{
		Env:    obj.Env,
		Logs:   obj.Logs,
		DB:     obj.DB,
		Repo:   r,
		OpenAI: obj.OpenAI,
	}
}
