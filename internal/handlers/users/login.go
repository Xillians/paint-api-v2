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

type LoginInput struct {
	GoogleUserId string `path:"id"`
}

type LoginOutputBody struct {
	Token     string `json:"token"`
	ExpiresAt string `json:"expires_at"`
}

type LoginOutput struct {
	Body LoginOutputBody
}

var loginOperation = huma.Operation{
	Method: http.MethodGet,
	Path:   "/login/{id}",
	Tags:   []string{"Users"},
}

func loginHandler(ctx context.Context, input *LoginInput) (*LoginOutput, error) {
	connection, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		slog.Error("could not retrieve db from context")
		return nil, huma.NewError(http.StatusInternalServerError, "failed to log in")
	}
	jwt, ok := ctx.Value("jwtKey").(jwt.JWTService)
	if !ok {
		slog.Error("could not retrieve jwt from context")
		return nil, huma.NewError(http.StatusInternalServerError, "failed to log in")
	}

	user, err := db.Users{}.GetUserByGoogleId(connection, input.GoogleUserId)
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
	return &LoginOutput{
		Body: LoginOutputBody{
			Token:     token,
			ExpiresAt: time.Now().Add(time.Hour * 7).String(),
		},
	}, nil
}
