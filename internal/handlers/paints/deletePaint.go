package paints

import (
	"context"
	"log/slog"
	"net/http"
	"paint-api/internal/db"
	"paint-api/internal/middleware"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

type DeletePaintInput struct {
	Id int `path:"id"`
}

type deletePaintOutput struct {
	Body string
}

var deleteOperation = huma.Operation{
	Method: http.MethodDelete,
	Path:   "/paints/{id}",
	Tags:   []string{"paints"},
}

func DeleteHandler(ctx context.Context, input *DeletePaintInput) (*deletePaintOutput, error) {
	connection, ok := ctx.Value(middleware.DbKey).(*gorm.DB)
	if !ok {
		slog.Error("Could not retrieve db from context")
		return nil, huma.NewError(http.StatusInternalServerError, "Failed to delete paint")
	}

	userRole, ok := ctx.Value(middleware.RoleKey).(string)
	if !ok {
		slog.Error("Could not retrieve user role from context")
		return nil, huma.NewError(http.StatusInternalServerError, "Failed to delete paint")
	}
	if userRole != "administrator" {
		return nil, huma.NewError(http.StatusForbidden, "You are not allowed to perform this action")
	}

	error := db.Paints{}.DeletePaint(connection, input.Id)
	if error != nil {
		slog.Error("Failed to delete paint", "error", error, "id", input.Id)
		return nil, huma.NewError(http.StatusInternalServerError, "Failed to delete paint")
	}

	return &deletePaintOutput{Body: "Paint deleted successfully"}, nil
}
