package middleware

import (
	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

type contextKey string

const DbKey contextKey = "db"

func UseDb(db *gorm.DB) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		ctx = huma.WithValue(ctx, DbKey, db)
		next(ctx)
	}
}
