package config

import (
	"log/slog"
	"os"

	"github.com/kelseyhightower/envconfig"
)

type DbConfig struct {
	DatabseUrl string `envconfig:"DATABASE_URL" required:"true"`
	AuthToken  string `envconfig:"AUTH_TOKEN" required:"true"`
}

type Config struct {
	DbConfig    DbConfig
	LogLevel    slog.Level `envconfig:"LOG_LEVEL" default:"info"`
	JwtSecret   string     `envconfig:"JWT_SECRET" required:"true"`
	HttpPort    int        `envconfig:"HTTP_PORT" default:"8080"`
	Environment string     `envconfig:"ENVIRONMENT" default:"development"`
}

func NewConfig() (*Config, error) {
	var c Config
	err := envconfig.Process("", &c)
	if err != nil {
		slog.Error("Failed to process environment variables", "error", err)
		return nil, err
	}

	if c.IsTest() {
		c.DbConfig.DatabseUrl = "file::memory:?cache=shared"
		c.DbConfig.AuthToken = "test-auth-token"
		c.JwtSecret = "test-jwt-secret"
	}

	c.LogLevel = ParseLogLevel(os.Getenv("LOG_LEVEL"))
	return &c, nil
}

func ParseLogLevel(level string) slog.Level {
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
func (c *Config) IsTest() bool {
	return c.Environment == "test"
}

func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}
