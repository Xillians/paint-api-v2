package paint_collection

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"paint-api/internal/db"
	"paint-api/internal/middleware"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

type UpdateCollectionEntryInputBody struct {
	Quantity int `json:"quantity"`
	PaintId  int `json:"paint_id"`
}

type UpdateCollectionEntryInput struct {
	Id   int `path:"id"`
	Body UpdateCollectionEntryInputBody
}

type UpdateCollectionEntryOutput struct {
	Body db.CollectionPaintDetails `json:"body"`
}

var updateOperation = huma.Operation{
	Method: http.MethodPut,
	Path:   "/collection/{id}",
	Tags:   []string{"collection"},
}

func UpdateHandler(ctx context.Context, input *UpdateCollectionEntryInput) (*UpdateCollectionEntryOutput, error) {
	connection, ok := ctx.Value(middleware.DbKey).(*gorm.DB)
	if !ok {
		slog.Error("Could not get database connection from context.")
		return nil, huma.NewError(http.StatusInternalServerError, "Failed to update entry.")
	}

	err := verifyCollectionOwnership(ctx, connection, input.Id)
	if err != nil {
		return nil, huma.NewError(http.StatusNotFound, "entry not found")
	}

	entry, err := db.CollectionPaintDetails{}.UpdateEntry(connection, db.UpdateCollectionEntryInput{
		ID:       input.Id,
		Quantity: input.Body.Quantity,
		PaintID:  input.Body.PaintId,
	})
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, huma.NewError(http.StatusNotFound, "entry not found")
		}
		slog.Error("An error occurred when updating collection entry.", "error", err)
		return nil, huma.NewError(http.StatusInternalServerError, "could not update entry")
	}

	return &UpdateCollectionEntryOutput{Body: *entry}, nil
}
