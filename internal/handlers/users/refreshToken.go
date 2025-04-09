package users

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"paint-api/internal/db"
	"paint-api/internal/jwt"
	"paint-api/internal/middleware"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

type RefreshTokenInput struct {
}

type refreshTokenOutputBody struct {
	Token     string `json:"token"`
	ExpiresAt string `json:"expires_at"`
}

type refreshTokenOutput struct {
	Body refreshTokenOutputBody
}

var refreshTokenOperation = huma.Operation{
	Method: http.MethodGet,
	Path:   "/refresh",
	Tags:   []string{"Users"},
}

func RefreshTokenHandler(ctx context.Context, input *RefreshTokenInput) (*refreshTokenOutput, error) {
	connection, ok := ctx.Value(middleware.DbKey).(*gorm.DB)
	if !ok {
		slog.Error("could not retrieve db from context")
		return nil, huma.NewError(http.StatusInternalServerError, "failed to refresh token")
	}
	userId, ok := ctx.Value(middleware.UserIdKey).(string)
	if !ok {
		slog.Error("could not retrieve userId from context")
		return nil, huma.NewError(http.StatusInternalServerError, "failed to refresh token")
	}
	jwt, ok := ctx.Value(middleware.JwtKey).(jwt.JWTService)
	if !ok {
		slog.Error("could not retrieve jwt from context")
		return nil, huma.NewError(http.StatusInternalServerError, "failed to refresh token")
	}

	user, err := db.Users{}.GetUserByGoogleId(connection, userId)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, huma.NewError(http.StatusNotFound, "User not found")
		}
		slog.Error("Error getting user", "error", err)
		return nil, huma.NewError(http.StatusInternalServerError, "Error getting user")
	}

	token, err := jwt.GenerateToken(user.GoogleUserId, user.Role)
	if err != nil {
		slog.Error("Error generating token", "error", err)
		return nil, huma.NewError(http.StatusInternalServerError, "Failed to refresh token")
	}
	return &refreshTokenOutput{
		Body: refreshTokenOutputBody{
			Token:     token,
			ExpiresAt: time.Now().Add(time.Hour * 7).String(),
		},
	}, nil
}
