package model

import "time"

type (
	EnvVar struct {
		AppApiHost   string
		AppApiPort   string
		LogEnv       string
		LogPath      string
		LogLevel     string
		TimeZone     *time.Location
		PSQLUrl		 string
		PSQLHost     string
		PSQLUsername string
		PSQLPassword string
		PSQLPort     string
		PSQLDB       string
		APIKey       string
		JWTSecret    string
		OpenAiApiKey string
		OpenAiModel  string
	}
)
