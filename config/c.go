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
	_ = godotenv.Load()

	loc, err := time.LoadLocation(os.Getenv("TIME_ZONE"))
	if err != nil || loc == nil {
		loc = time.UTC
	}

	env := &model.EnvVar{
		AppApiHost: os.Getenv("APP_API_HOST"),
		AppApiPort: os.Getenv("APP_API_PORT"),
		LogEnv:     os.Getenv("LOG_ENV"),
		LogPath:    os.Getenv("LOG_PATH"),
		LogLevel:   os.Getenv("LOG_LEVEL"),
		TimeZone:   loc,

		// DB
		PSQLUrl:      os.Getenv("PSQL_URL"),
		PSQLHost:     os.Getenv("PSQL_HOST"),
		PSQLUsername: os.Getenv("PSQL_USERNAME"),
		PSQLPassword: os.Getenv("PSQL_PASSWORD"),
		PSQLPort:     os.Getenv("PSQL_PORT"),
		PSQLDB:       os.Getenv("PSQL_DB"),

		// App/API
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
	// Prefer PSQL_URL when present
	dsn := strings.TrimSpace(env.PSQLUrl)

	if dsn == "" {
		// Build DSN from individual parts
		host := strings.TrimSpace(env.PSQLHost)
		user := strings.TrimSpace(env.PSQLUsername)
		pass := strings.TrimSpace(env.PSQLPassword)
		port := strings.TrimSpace(env.PSQLPort)
		db := strings.TrimSpace(env.PSQLDB)

		// Sensible defaults
		if port == "" {
			port = "5432"
		}
		sslmode := strings.TrimSpace(os.Getenv("PSQL_SSLMODE"))
		if sslmode == "" {
			// For local/dev commonly "disable"; override via PSQL_SSLMODE if needed
			sslmode = "disable"
		}
		tz := "UTC"
		if env.TimeZone != nil {
			tz = env.TimeZone.String()
		}

		// Basic validation of required parts
		if host == "" || user == "" || db == "" {
			log.Fatal("database configuration incomplete: require PSQL_HOST, PSQL_USERNAME, and PSQL_DB when PSQL_URL is not set")
		}

		// Final DSN
		dsn = fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
			host, user, pass, db, port, sslmode, tz,
		)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		l.Errorf("[x] Error establishing database connection: %v", err)
		panic(err)
	}

	// Enable useful extensions (idempotent)
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "pgcrypto"`).Error; err != nil {
		log.Fatalf("enable pgcrypto: %v", err)
	}
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS unaccent;`).Error; err != nil {
		log.Fatalf("enable unaccent : %v", err)
	}

	l.Info("[v] Database GORM successfully established")

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db
}

func InitOpenAIService(l *logrus.Logger, env *model.EnvVar) *openai.Client {
	key := strings.TrimSpace(env.OpenAiApiKey)
	if key == "" {
		l.Fatal("OPENAI_API_KEY is empty")
	}

	base := strings.TrimSpace(os.Getenv("OPENAI_BASE_URL"))
	if base == "" {
		base = "https://api.openai.com/v1"
	}

	client := openai.NewClient(
		option.WithAPIKey(key),
		option.WithBaseURL(base),
	)

	// Test connectivity (non-fatal if it fails)
	if _, err := client.Models.List(context.Background()); err != nil {
		l.WithError(err).Error("failed to connect to OpenAI API")
	} else {
		l.Infof("OpenAI connection is OK, base=%s", base)
	}
	return &client
}
