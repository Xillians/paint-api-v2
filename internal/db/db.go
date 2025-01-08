package db

import (
	"fmt"
	"log/slog"
	"paint-api/internal/config"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

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
