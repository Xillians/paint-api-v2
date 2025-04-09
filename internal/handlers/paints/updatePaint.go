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

type UpdatePaintInputBody struct {
	Name string `json:"name"`
}
type UpdatePaintInput struct {
	Id   int `path:"id"`
	Body UpdatePaintInputBody
}

type updatePaintOutput struct {
	Body db.Paints
}

var updateOperation = huma.Operation{
	Method: http.MethodPut,
	Path:   "/paints/{id}",
	Tags:   []string{"paints"},
}

func UpdateHandler(ctx context.Context, input *UpdatePaintInput) (*updatePaintOutput, error) {
	userRole, ok := ctx.Value(middleware.RoleKey).(string)
	if !ok {
		slog.Error("Could not retrieve user role from context")
		return nil, huma.NewError(http.StatusInternalServerError, "failed to update paint")
	}
	if userRole != "administrator" {
		return nil, huma.NewError(http.StatusForbidden, "You are not allowed to perform this action")
	}

	connection, ok := ctx.Value(middleware.DbKey).(*gorm.DB)
	if !ok {
		slog.Error("Could not retrieve db from context")
		return nil, huma.NewError(http.StatusInternalServerError, "failed to update paint")
	}

	dbInput := &db.UpdatePaintInput{
		Name: input.Body.Name,
	}

	paint, err := db.Paints{}.UpdatePaint(connection, input.Id, dbInput)
	if err != nil {
		if err == db.ErrRecordNotFound {
			return nil, huma.NewError(http.StatusNotFound, "paint not found")
		}
		slog.Error("Failed to update paint", "error", err, "id", input.Id)
		return nil, huma.NewError(http.StatusInternalServerError, "failed to update paint")
	}

	return &updatePaintOutput{Body: *paint}, nil
}
