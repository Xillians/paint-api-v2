package middleware

import (
	"net/http"
	"paint-api/internal/jwt"
	"regexp"
	"strings"

	j "github.com/golang-jwt/jwt/v5"

	"github.com/danielgtaylor/huma/v2"
)

const RoleKey contextKey = "role"
const UserIdKey contextKey = "userId"
const JwtKey contextKey = "jwtKey"

func UseJwt(jwtService jwt.JWTService) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		ctx = huma.WithValue(ctx, JwtKey, jwtService)
		next(ctx)
	}
}

func AuthenticateRequests(api huma.API, jwtService jwt.JWTService) func(huma.Context, func(huma.Context)) {
	return func(ctx huma.Context, next func(huma.Context)) {
		if isLoginOrRegisterPath(ctx) {
			next(ctx)
			return
		}

		authHeader := ctx.Header("Authorization")
		jwtToken := strings.Split(authHeader, " ")[1]
		token, err := jwtService.VerifyToken(jwtToken)
		if err != nil {
			_ = huma.WriteErr(api, ctx, http.StatusUnauthorized, "Unauthorized", err)
			return
		}

		userId := token.Claims.(j.MapClaims)["user_id"]
		if userId == nil {
			_ = huma.WriteErr(api, ctx, 401, "Unauthorized: User ID not found in token", nil)
			return
		}
		role := token.Claims.(j.MapClaims)["role"]
		if role == nil {
			_ = huma.WriteErr(api, ctx, 401, "Unauthorized: Role not found in token", nil)
			return
		}

		ctx = huma.WithValue(ctx, UserIdKey, userId)
		ctx = huma.WithValue(ctx, RoleKey, role)

		next(ctx)
	}
}

func isLoginOrRegisterPath(ctx huma.Context) bool {
	path := ctx.URL().Path
	loginPathRegex := regexp.MustCompile(`^/login/\d+$`)
	return loginPathRegex.MatchString(path) || path == "/register"
}
