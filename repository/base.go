package repository

import (
	"cms-octo-chat-api/model"

	"cloud.google.com/go/pubsub"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type (
	BaseRepository struct {
		Env    *model.EnvVar
		Logs   *logrus.Logger
		DB     *gorm.DB
		PubSub *pubsub.Client
	}
)

func NewBaseRepository(obj BaseRepository) *BaseRepository {

	return &BaseRepository{
		Env:    obj.Env,
		Logs:   obj.Logs,
		DB:     obj.DB,
		PubSub: obj.PubSub,
	}
}
