package users

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"paint-api/internal/db"
	"paint-api/internal/middleware"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

type ForgetUserInput struct {
}

type forgetUserOutput struct {
	Body string
}

var forgetOperation = huma.Operation{
	Method: http.MethodDelete,
	Path:   "/forget",
	Tags:   []string{"Users"},
}

func ForgetHandler(ctx context.Context, input *ForgetUserInput) (*forgetUserOutput, error) {
	connection, ok := ctx.Value(middleware.DbKey).(*gorm.DB)
	if !ok {
		slog.Error("could not retrieve db from context")
		return nil, huma.NewError(http.StatusInternalServerError, "failed to delete user")
	}
	userId, ok := ctx.Value(middleware.UserIdKey).(string)
	if !ok {
		slog.Error("could not retrieve userId from context")
		return nil, huma.NewError(http.StatusInternalServerError, "failed to delete user")
	}

	err := db.Users{}.DeleteUserByGoogleId(connection, userId)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, huma.NewError(http.StatusNotFound, "User not found")
		}

		slog.Error("Error deleting user", "error", err)
		return nil, huma.NewError(http.StatusInternalServerError, "Error deleting user")
	}

	return &forgetUserOutput{Body: "User deleted successfully"}, nil
}
