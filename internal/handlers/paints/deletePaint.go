package paints

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"paint-api/internal/db"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

type deletePaintInput struct {
	Id int `path:"id"`
}

type deletePaintOutput struct {
	Body string
}

var DeleteOperation = huma.Operation{
	Method: http.MethodDelete,
	Path:   "/paints/{id}",
	Tags:   []string{"paints"},
}

func DeleteHandler(ctx context.Context, input *deletePaintInput) (*deletePaintOutput, error) {
	connection, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return nil, errors.New("could not retrieve db from context")
	}

	userRole := ctx.Value("role").(string)
	if userRole != "administrator" {
		return nil, huma.NewError(403, "You are not allowed to perform this action")
	}

	error := db.Paints{}.DeletePaint(connection, input.Id)
	if error != nil {
		slog.Error("Failed to delete paint", "error", error, "id", input.Id)
		return nil, huma.NewError(http.StatusInternalServerError, "Failed to delete paint")
	}

	return &deletePaintOutput{Body: "Paint deleted successfully"}, nil
}
