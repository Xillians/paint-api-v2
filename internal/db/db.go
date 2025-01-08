package db

import (
	"database/sql"
	"fmt"
	"log/slog"
	"paint-api/internal/config"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func New(cfg *config.DbConfig) (*sql.DB, error) {
	url := fmt.Sprintf("%s?authToken=%s", cfg.DatabseUrl, cfg.AuthToken)
	slog.Debug("Connecting to database", slog.Any("url", url))

	db, err := sql.Open("libsql", url)
	if err != nil {
		return nil, err
	}
	return db, nil
}
