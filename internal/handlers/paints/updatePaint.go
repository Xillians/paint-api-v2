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

type updatePaintInputBody struct {
	Name string `json:"name"`
}
type updatePaintInput struct {
	Id   int `path:"id"`
	Body updatePaintInputBody
}

type updatePaintOutput struct {
	Body db.CollectionPaintDetails
}

var UpdateOperation = huma.Operation{
	Method: http.MethodPut,
	Path:   "/paints/{id}",
	Tags:   []string{"paints"},
}

func UpdateHandler(ctx context.Context, input *updatePaintInput) (*updatePaintOutput, error) {
	out := updatePaintOutput{
		Body: db.CollectionPaintDetails{},
	}

	userRole := ctx.Value("role").(string)
	if userRole != "administrator" {
		return nil, huma.NewError(403, "You are not allowed to perform this action")
	}

	connection, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		slog.Error("Could not retrieve db from context")
		return nil, huma.NewError(http.StatusInternalServerError, "failed to update paint")
	}
	var paint db.Paints
	if err := connection.First(&paint, input.Id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, huma.NewError(http.StatusNotFound, "paint not found")
		}
		slog.Error("An error occurred when fetching paint.", "error", err)
		return nil, huma.NewError(http.StatusInternalServerError, "could not fetch paint")
	}

	paint.Name = input.Body.Name
	connection.Save(&paint)

	err := connection.Preload("Brand").Find(&out.Body).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, huma.NewError(http.StatusNotFound, "paint not found")
		}
		slog.Error("An error occurred when fetching paint.", "error", err)
		return nil, huma.NewError(http.StatusInternalServerError, "could not fetch paint")
	}

	return &out, nil
}
