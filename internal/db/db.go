package db

import (
	"errors"
	"fmt"
	"log/slog"
	"paint-api/internal/config"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Add error type: Row not found for local use
var ErrRecordNotFound = gorm.ErrRecordNotFound
var ErrRecordExists = errors.New("Record already exists")

func New(cfg *config.DbConfig) (*gorm.DB, error) {
	url := fmt.Sprintf("%s?authToken=%s", cfg.DatabseUrl, cfg.AuthToken)
	slog.Debug("Connecting to database", slog.Any("url", url))

	db, err := gorm.Open(sqlite.New(sqlite.Config{
		DriverName: "libsql",
		DSN:        url,
	}), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	return db, nil
}
