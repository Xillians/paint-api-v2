package users

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"paint-api/internal/db"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

type forgetUserInput struct {
}

type forgetUserOutput struct {
	Body string
}

var ForgetOperation = huma.Operation{
	Method: http.MethodDelete,
	Path:   "/forget",
	Tags:   []string{"Users"},
}

func ForgetHandler(ctx context.Context, input *forgetUserInput) (*forgetUserOutput, error) {
	connection, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		slog.Error("could not retrieve db from context")
		return nil, huma.NewError(http.StatusInternalServerError, "failed to delete user")
	}
	userId, ok := ctx.Value("userId").(string)
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
