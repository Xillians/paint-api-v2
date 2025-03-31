package middleware

import (
	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

func UseDb(db *gorm.DB) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		ctx = huma.WithValue(ctx, "db", db)
		next(ctx)
	}
}
