package paints

import (
	"context"
	"log/slog"
	"net/http"
	"paint-api/internal/db"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

type listPaintInput struct {
}

type listPaintOutputBody struct {
	Paints []db.Paints `json:"paints"`
}

type listPaintOutput struct {
	Body listPaintOutputBody `json:"body"`
}

var listOperation = huma.Operation{
	Method: http.MethodGet,
	Path:   "/paints",
	Tags:   []string{"paints"},
}

func listHandler(ctx context.Context, input *listPaintInput) (*listPaintOutput, error) {
	connection, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		slog.Error("Could not retrieve db from context")
		return nil, huma.NewError(http.StatusInternalServerError, "failed to list paints")
	}

	paints, err := db.Paints{}.ListPaints(connection)
	if err != nil {
		slog.Error("An error occurred when fetching paints.", "error", err)
		return nil, huma.NewError(http.StatusInternalServerError, "could not fetch paints")
	}

	return &listPaintOutput{Body: listPaintOutputBody{Paints: paints}}, nil
}
