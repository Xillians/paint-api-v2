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

var ListPaintCollectionOperation = huma.Operation{
	Method: http.MethodGet,
	Path:   "/collection",
	Tags:   []string{"collection"},
}

func ListPaintCollectionHandler(ctx context.Context, input *listPaintCollectionInput) (*listPaintCollectionOutput, error) {
	out := listPaintCollectionOutput{
		Body: listPaintCollectionOutputBody{
			Collection: []db.CollectionPaintDetails{},
		},
	}
	userId, ok := ctx.Value("userId").(string)
	if !ok {
		return nil, errors.New("could not retrieve user_id from context")
	}
	connection, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return nil, errors.New("could not retrieve db from context")
	}

	err := connection.
		Preload("Paint").
		Preload("Paint.Brand").
		Joins("JOIN users ON users.id = paint_collections.user_id").
		Where("users.google_user_id = ?", userId).
		Find(&out.Body.Collection).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, huma.NewError(http.StatusNotFound, "user not found")
		}
		slog.Error("An error occurred when fetching collection.", "error", err)
		return nil, huma.NewError(http.StatusInternalServerError, "could not fetch collection")
	}

	return &out, nil
}
