package middleware

import (
	"database/sql"

	"github.com/danielgtaylor/huma/v2"
)

type key int

const dbKey key = 0

func UseDb(db *sql.DB) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		ctx = huma.WithValue(ctx, "db", db)
		next(ctx)
	}
}
