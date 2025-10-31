package config

import (
	"cms-octo-chat-api/helper"
	"cms-octo-chat-api/model"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	openai "github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func LoadEnv() *model.EnvVar {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	loc, _ := time.LoadLocation(os.Getenv("TIME_ZONE"))
	env := &model.EnvVar{
		AppApiHost:   os.Getenv("APP_API_HOST"),
		AppApiPort:   os.Getenv("APP_API_PORT"),
		LogEnv:       os.Getenv("LOG_ENV"),
		LogPath:      os.Getenv("LOG_PATH"),
		LogLevel:     os.Getenv("LOG_LEVEL"),
		TimeZone:     loc,
		PSQLHost:     os.Getenv("PSQL_HOST"),
		PSQLUsername: os.Getenv("PSQL_USERNAME"),
		PSQLPassword: os.Getenv("PSQL_PASSWORD"),
		PSQLPort:     os.Getenv("PSQL_PORT"),
		PSQLDB:       os.Getenv("PSQL_DB"),
		APIKey:       os.Getenv("API_KEY"),
		JWTSecret:    os.Getenv("JWT_SECRET"),
		OpenAiModel:  os.Getenv("OPENAI_MODEL"),
		OpenAiApiKey: os.Getenv("OPENAI_API_KEY"),
	}
	return env
}

func Initiate(logname string) *model.Resources {
	env := LoadEnv()

	l := helper.MakeLogger(
		helper.Setup{Env: env.LogEnv, Logname: env.LogPath + "/" + logname, Display: true, Level: env.LogLevel})
	l.Info(fmt.Sprintf("Environment Variable Loaded : %#v\n", env))

	return &model.Resources{
		Env:    env,
		Logs:   l,
		DB:     initGormPsql(l, env),
		OpenAI: InitOpenAIService(l, env),
	}
}

func initGormPsql(l *logrus.Logger, env *model.EnvVar) *gorm.DB {

	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  "host=" + env.PSQLHost + " user=" + env.PSQLUsername + " password=" + env.PSQLPassword + " dbname=" + env.PSQLDB + " port=" + env.PSQLPort + " sslmode=disable TimeZone=" + env.TimeZone.String(),
		PreferSimpleProtocol: true,
	}), &gorm.Config{})

	if err != nil {

		// panic the function then hard exit
		l.Info(fmt.Sprintf("[x] An Error occured when establishing of the database : %#v\n", err))

		panic(err)

	} else {

		if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "pgcrypto"`).Error; err != nil {
			log.Fatalf("enable pgcrypto: %v", err)
		}
		if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS unaccent;`).Error; err != nil {
			log.Fatalf("enable unaccent : %v", err)
		}

		l.Info("[v] Database GORM successful established\n")

		sqlDB, _ := db.DB()

		// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
		sqlDB.SetMaxIdleConns(10)

		// SetMaxOpenConns sets the maximum number of open connections to the database.
		sqlDB.SetMaxOpenConns(100)

		// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
		sqlDB.SetConnMaxLifetime(time.Hour)
	}

	return db
}

func InitOpenAIService(l *logrus.Logger, env *model.EnvVar) *openai.Client {
	key := env.OpenAiApiKey
	if key == "" {
		l.Fatal("OPENAI_API_KEY is empty")
	}

	base := os.Getenv("OPENAI_BASE_URL")
	if strings.TrimSpace(base) == "" {
		base = "https://api.openai.com/v1" // default resmi
	}

	client := openai.NewClient(
		option.WithAPIKey(key),
		option.WithBaseURL(base),
	)

	// test connectivity
	_, err := client.Models.List(context.Background())
	if err != nil {
		l.WithError(err).Error("failed to connect to OpenAI API")
	} else {
		l.Infof("OpenAI connection is OK, base=%s", base)
	}
	return &client
}
