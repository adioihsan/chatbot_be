package model

import (
	openai "github.com/openai/openai-go/v3"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type (
	GlobalResponse struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    any    `json:"data"`
	}

	Resources struct {
		Env    *EnvVar
		Logs   *logrus.Logger
		DB     *gorm.DB
		OpenAI *openai.Client
	}

	ErrorValidationResponse struct {
		Code    int               `json:"code"`
		Message string            `json:"message"`
		Errors  map[string]string `json:"errors,omitempty"`
	}
)
