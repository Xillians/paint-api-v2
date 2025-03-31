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

type listPaintInput struct {
}

type listPaintOutputBody struct {
	Paints []db.Paints `json:"paints"`
}

type listPaintOutput struct {
	Body listPaintOutputBody `json:"body"`
}

var ListOperation = huma.Operation{
	Method: http.MethodGet,
	Path:   "/paints",
	Tags:   []string{"paints"},
}

func ListHandler(ctx context.Context, input *listPaintInput) (*listPaintOutput, error) {
	connection, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return nil, errors.New("could not retrieve db from context")
	}

	paints, err := db.Paints{}.ListPaints(connection)
	if err != nil {
		slog.Error("An error occurred when fetching paints.", "error", err)
		return nil, huma.NewError(http.StatusInternalServerError, "could not fetch paints")
	}

	return &listPaintOutput{Body: listPaintOutputBody{Paints: paints}}, nil
}
