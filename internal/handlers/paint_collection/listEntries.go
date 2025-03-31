package paint_collection

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"paint-api/internal/db"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

type listPaintCollectionInput struct {
}

type listPaintCollectionOutputBody struct {
	Collection []db.CollectionPaintDetails `json:"collection"`
}

type listPaintCollectionOutput struct {
	Body listPaintCollectionOutputBody `json:"body"`
}

var listOperation = huma.Operation{
	Method: http.MethodGet,
	Path:   "/collection",
	Tags:   []string{"collection"},
}

func listHandler(ctx context.Context, input *listPaintCollectionInput) (*listPaintCollectionOutput, error) {
	userId, ok := ctx.Value("userId").(string)
	if !ok {
		slog.Error("could not retrieve userId from context")
		return nil, huma.NewError(http.StatusInternalServerError, "failed to list collection entries")
	}
	connection, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		slog.Error("could not retrieve db from context")
		return nil, huma.NewError(http.StatusInternalServerError, "failed to list collection entries")
	}

	entries, err := db.CollectionPaintDetails{}.ListEntries(connection, userId)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return &listPaintCollectionOutput{
				Body: listPaintCollectionOutputBody{
					Collection: []db.CollectionPaintDetails{},
				},
			}, nil
		}
		slog.Error("failed to list collection entries", "error", err)
		return nil, huma.NewError(http.StatusInternalServerError, "failed to list collection entries")
	}

	return &listPaintCollectionOutput{
		Body: listPaintCollectionOutputBody{
			Collection: entries,
		},
	}, nil
}
