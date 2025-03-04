package config

import (
	"log"
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

func NewConfig() *Config {
	var c Config
	err := envconfig.Process("", &c)
	if err != nil {
		log.Fatal(err)
	}
	c.LogLevel = parseLogLevel(os.Getenv("LOG_LEVEL"))
	slog.Info("Config", "config", c)
	return &c
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
