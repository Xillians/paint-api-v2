package middleware

import (
	"paint-api/internal/jwt"
	"regexp"
	"strings"

	"github.com/danielgtaylor/huma/v2"
)

const jwtKey key = 1

func UseJwt(jwtService jwt.JWTService) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		ctx = huma.WithValue(ctx, "jwtKey", jwtService)
		next(ctx)
	}
}

func AuthenticateRequests(api huma.API, jwtService jwt.JWTService) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		loginPathRegex := regexp.MustCompile(`^/login/\d+$`)
		if loginPathRegex.MatchString(ctx.URL().Path) || ctx.URL().Path == "/register" {
			next(ctx)
			return
		}

		authHeader := ctx.Header("Authorization")
		jwt := strings.Split(authHeader, " ")[1]
		_, err := jwtService.VerifyToken(jwt)
		if err != nil {
			huma.WriteErr(api, ctx, 401, "Unauthorized", err)
			return
		}

		next(ctx)
	}
}
