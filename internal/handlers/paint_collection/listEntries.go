package paint_collection

import (
	"context"
	"log/slog"
	"net/http"
	"paint-api/internal/db"
	"paint-api/internal/middleware"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

type ListPaintCollectionInput struct {
}

type ListPaintCollectionOutputBody struct {
	Collection []db.CollectionPaintDetails `json:"collection"`
}

type ListPaintCollectionOutput struct {
	Body ListPaintCollectionOutputBody `json:"body"`
}

var listOperation = huma.Operation{
	Method: http.MethodGet,
	Path:   "/collection",
	Tags:   []string{"collection"},
}

func ListHandler(ctx context.Context, input *ListPaintCollectionInput) (*ListPaintCollectionOutput, error) {
	userId, ok := ctx.Value(middleware.UserIdKey).(string)
	if !ok {
		slog.Error("could not retrieve userId from context")
		return nil, huma.NewError(http.StatusInternalServerError, "failed to list collection entries")
	}
	connection, ok := ctx.Value(middleware.DbKey).(*gorm.DB)
	if !ok {
		slog.Error("could not retrieve db from context")
		return nil, huma.NewError(http.StatusInternalServerError, "failed to list collection entries")
	}

	entries, err := db.CollectionPaintDetails{}.ListEntries(connection, userId)
	if err != nil {
		slog.Error("failed to list collection entries", "error", err)
		return nil, huma.NewError(http.StatusInternalServerError, "failed to list collection entries")
	}

	return &ListPaintCollectionOutput{
		Body: ListPaintCollectionOutputBody{
			Collection: entries,
		},
	}, nil
}
