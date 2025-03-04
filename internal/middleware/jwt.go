package middleware

import (
	"paint-api/internal/jwt"

	"github.com/danielgtaylor/huma/v2"
)

const jwtKey key = 1

func UseJwt(jwtService jwt.JWTService) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		ctx = huma.WithValue(ctx, "jwtKey", jwtService)
		next(ctx)
	}
}
