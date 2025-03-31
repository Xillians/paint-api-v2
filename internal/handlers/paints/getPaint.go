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

type getPaintsInput struct {
	Id int `path:"id"`
}

type getPaintOutput struct {
	Body db.Paints `json:"body"`
}

var getOperation = huma.Operation{
	Method: http.MethodGet,
	Path:   "/paints/{id}",
	Tags:   []string{"paints"},
}

func getHandler(ctx context.Context, input *getPaintsInput) (*getPaintOutput, error) {
	connection, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		slog.Error("Could not retrieve db from context")
		return nil, huma.NewError(http.StatusInternalServerError, "could not fetch paint")
	}

	paints, err := db.Paints{}.GetPaint(connection, input.Id)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, huma.NewError(http.StatusNotFound, "paint not found")
		}
		slog.Error("An error occurred when fetching paint.", "error", err)
		return nil, huma.NewError(http.StatusInternalServerError, "could not fetch paint")
	}

	return &getPaintOutput{Body: *paints}, nil
}
