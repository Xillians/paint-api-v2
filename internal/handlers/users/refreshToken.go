package users

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"paint-api/internal/db"
	"paint-api/internal/jwt"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

type refreshTokenInput struct {
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

func refreshTokenHandler(ctx context.Context, input *refreshTokenInput) (*refreshTokenOutput, error) {
	connection, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		slog.Error("could not retrieve db from context")
		return nil, huma.NewError(http.StatusInternalServerError, "failed to refresh token")
	}
	userId, ok := ctx.Value("userId").(string)
	if !ok {
		slog.Error("could not retrieve userId from context")
		return nil, huma.NewError(http.StatusInternalServerError, "failed to refresh token")
	}
	jwt, ok := ctx.Value("jwtKey").(jwt.JWTService)
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
		return nil, err
	}
	return &refreshTokenOutput{
		Body: refreshTokenOutputBody{
			Token:     token,
			ExpiresAt: time.Now().Add(time.Hour * 7).String(),
		},
	}, nil
}
