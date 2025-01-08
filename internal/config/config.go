package config

import (
	"log"
	"log/slog"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type DbConfig struct {
	DatabseUrl string `json:"database_url"`
	AuthToken  string `json:"auth_token"`
}

type Config struct {
	DbConfig    DbConfig
	LogLevel    slog.Level `json:"log_level" default:"info"`
	HttpPort    int        `json:"http_port" default:"8080"`
	Environment string     `json:"environment" default:"development"`
}

func NewConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	port, err := strconv.Atoi(os.Getenv("HTTP_PORT"))
	if err != nil {
		log.Fatalf("Error parsing HTTP_PORT")
	}

	return &Config{
		DbConfig: DbConfig{
			DatabseUrl: os.Getenv("DATABASE_URL"),
			AuthToken:  os.Getenv("AUTH_TOKEN"),
		},
		LogLevel:    parseLogLevel(os.Getenv("LOG_LEVEL")),
		HttpPort:    port,
		Environment: os.Getenv("ENVIRONMENT"),
	}
}

func parseLogLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func (c *Config) GetLogLevel() slog.Level {
	return c.LogLevel
}

func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}
