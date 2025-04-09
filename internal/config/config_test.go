package config_test

import (
	"log/slog"
	"os"
	"paint-api/internal/config"
	"testing"
)

func TestNewConfig(t *testing.T) {
	dbUrl := "file::memory:?cache=shared"
	authToken := "test-auth-token"
	logLevel := "info"
	jwtSecret := "some_secret"
	environment := "development"
	os.Setenv("DATABASE_URL", dbUrl)
	os.Setenv("AUTH_TOKEN", authToken)
	os.Setenv("LOG_LEVEL", logLevel)
	os.Setenv("JWT_SECRET", jwtSecret)
	os.Setenv("ENVIRONMENT", environment)
	t.Run("Create default config", func(t *testing.T) {
		cfg, err := config.NewConfig()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if cfg.DbConfig.DatabaseUrl != dbUrl {
			t.Errorf("Expected %s, got %s", dbUrl, cfg.DbConfig.DatabaseUrl)
		}
		if cfg.DbConfig.AuthToken != authToken {
			t.Errorf("Expected %s, got %s", authToken, cfg.DbConfig.AuthToken)
		}
		if cfg.LogLevel != config.ParseLogLevel(logLevel) {
			t.Errorf("Expected %s, got %s", logLevel, cfg.LogLevel)
		}
		if cfg.JwtSecret != jwtSecret {
			t.Errorf("Expected %s, got %s", jwtSecret, cfg.JwtSecret)
		}
		if cfg.Environment != environment {
			t.Errorf("Expected %s, got %s", environment, cfg.Environment)
		}
	})
	t.Run("Get testing config", func(t *testing.T) {
		os.Setenv("ENVIRONMENT", "test")
		cfg, err := config.NewConfig()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if cfg.DbConfig.DatabaseUrl != "file::memory:?cache=shared" {
			t.Errorf("Expected %s, got %s", "file::memory:?cache=shared", cfg.DbConfig.DatabaseUrl)
		}
		if cfg.DbConfig.AuthToken != "test-auth-token" {
			t.Errorf("Expected %s, got %s", "test-auth-token", cfg.DbConfig.AuthToken)
		}
		if cfg.JwtSecret != "test-jwt-secret" {
			t.Errorf("Expected test-jwt-secret, got %s", cfg.JwtSecret)
		}
	})
	t.Run("missing required env variables", func(t *testing.T) {
		os.Unsetenv("DATABASE_URL")
		os.Unsetenv("AUTH_TOKEN")
		os.Unsetenv("JWT_SECRET")
		_, err := config.NewConfig()
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})
	t.Cleanup(func() {
		os.Unsetenv("DATABASE_URL")
		os.Unsetenv("AUTH_TOKEN")
		os.Unsetenv("LOG_LEVEL")
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("HTTP_PORT")
		os.Unsetenv("ENVIRONMENT")
	})
}

func TestParseLogLevel(t *testing.T) {
	dbUrl := "file::memory:?cache=shared"
	authToken := "test-auth-token"
	logLevel := "info"
	jwtSecret := "some_secret"
	environment := "development"
	os.Setenv("DATABASE_URL", dbUrl)
	os.Setenv("AUTH_TOKEN", authToken)
	os.Setenv("LOG_LEVEL", logLevel)
	os.Setenv("JWT_SECRET", jwtSecret)
	os.Setenv("ENVIRONMENT", environment)

	logLevels := map[string]slog.Level{
		"debug": slog.LevelDebug,
		"info":  slog.LevelInfo,
		"warn":  slog.LevelWarn,
		"error": slog.LevelError,
	}
	for level, expected := range logLevels {
		t.Run(level, func(t *testing.T) {
			os.Setenv("LOG_LEVEL", level)
			cfg, err := config.NewConfig()
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
			if cfg.LogLevel != expected {
				t.Errorf("Expected %s, got %s", expected, cfg.LogLevel)
			}
			level := cfg.GetLogLevel()
			if level != expected {
				t.Errorf("Expected %s, got %s", expected, level)
			}
		})
	}
	t.Run("default", func(t *testing.T) {
		os.Unsetenv("LOG_LEVEL")
		cfg, err := config.NewConfig()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if cfg.LogLevel != slog.LevelInfo {
			t.Errorf("Expected %s, got %s", slog.LevelInfo, cfg.LogLevel)
		}
		level := cfg.GetLogLevel()
		if level != slog.LevelInfo {
			t.Errorf("Expected %s, got %s", slog.LevelInfo, level)
		}
	})
	t.Cleanup(func() {
		os.Unsetenv("DATABASE_URL")
		os.Unsetenv("AUTH_TOKEN")
		os.Unsetenv("LOG_LEVEL")
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("ENVIRONMENT")
	})
}

func TestEnvironmentLevels(t *testing.T) {
	dbUrl := "file::memory:?cache=shared"
	authToken := "test-auth-token"
	logLevel := "info"
	jwtSecret := "some_secret"
	os.Setenv("DATABASE_URL", dbUrl)
	os.Setenv("AUTH_TOKEN", authToken)
	os.Setenv("LOG_LEVEL", logLevel)
	os.Setenv("JWT_SECRET", jwtSecret)
	t.Run("Development", func(t *testing.T) {
		os.Setenv("ENVIRONMENT", "development")
		cfg, err := config.NewConfig()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if !cfg.IsDevelopment() {
			t.Errorf("Expected true, got false")
		}
	})
	t.Run("Test", func(t *testing.T) {
		os.Setenv("ENVIRONMENT", "test")
		cfg, err := config.NewConfig()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if !cfg.IsTest() {
			t.Errorf("Expected true, got false")
		}
	})
	t.Run("Production", func(t *testing.T) {
		os.Setenv("ENVIRONMENT", "production")
		cfg, err := config.NewConfig()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if !cfg.IsProduction() {
			t.Errorf("Expected true, got false")
		}
	})
	t.Run("Default", func(t *testing.T) {
		os.Unsetenv("ENVIRONMENT")
		cfg, err := config.NewConfig()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if !cfg.IsDevelopment() {
			t.Errorf("Expected true, got false")
		}
	})
	t.Cleanup(func() {
		os.Unsetenv("DATABASE_URL")
		os.Unsetenv("AUTH_TOKEN")
		os.Unsetenv("LOG_LEVEL")
		os.Unsetenv("JWT_SECRET")
		os.Unsetenv("ENVIRONMENT")
	})
}
