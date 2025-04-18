package users

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"paint-api/internal/db"
	"paint-api/internal/middleware"
	"regexp"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

type RegisterUserInput struct {
	Body db.RegisterUserInput `json:"body"`
}

type RegisterUserOutput struct {
	Body db.Users
}

var registerOperation = huma.Operation{
	Method: http.MethodPost,
	Path:   "/register",
	Tags:   []string{"Users"},
}

func RegisterHandler(ctx context.Context, input *RegisterUserInput) (*RegisterUserOutput, error) {
	connection, ok := ctx.Value(middleware.DbKey).(*gorm.DB)
	if !ok {
		slog.Error("could not retrieve db from context")
		return nil, huma.NewError(http.StatusInternalServerError, "failed to register user")
	}
	if !validateEmail(input.Body.Email) {
		return nil, huma.NewError(http.StatusBadRequest, "invalid email")
	}

	user, err := db.Users{}.RegisterUser(connection, input.Body, "user")
	if err != nil {
		if errors.Is(err, db.ErrRecordExists) {
			return nil, huma.NewError(http.StatusConflict, "User already exists")
		}
		slog.Error("Error registering user", "error", err)
		return nil, huma.NewError(http.StatusInternalServerError, "Error registering user")
	}

	return &RegisterUserOutput{Body: *user}, nil
}

func validateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(email)
}
