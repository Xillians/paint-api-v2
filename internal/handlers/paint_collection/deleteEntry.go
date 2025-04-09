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

type DeleteCollectionEntryInput struct {
	Id int `path:"id"`
}

type DeleteCollectionEntryOutput struct {
	Body string
}

var deleteOperation = huma.Operation{
	Method: http.MethodDelete,
	Path:   "/collection/{id}",
	Tags:   []string{"collection"},
}

// DeleteHandler deletes a paint from the collection
//
// Deletes a paint from the collection. Requires authentication, and the paint must belong to the user.
func DeleteHandler(ctx context.Context, input *DeleteCollectionEntryInput) (*DeleteCollectionEntryOutput, error) {
	connection, ok := ctx.Value(middleware.DbKey).(*gorm.DB)
	if !ok {
		slog.Error("could not retrieve db from context")
		return nil, huma.NewError(http.StatusInternalServerError, "failed to delete entry")
	}

	err := verifyCollectionOwnership(ctx, connection, input.Id)
	if err != nil {
		slog.Error("could not verify collection ownership", "error", err)
		return nil, huma.NewError(http.StatusNotFound, "Entry not found")
	}

	err = db.CollectionPaintDetails{}.DeleteEntry(connection, input.Id)
	if err != nil {
		slog.Error("Error deleting paint", "error", err)
		return nil, huma.NewError(http.StatusInternalServerError, "Error deleting entry")
	}
	return &DeleteCollectionEntryOutput{Body: "Entry deleted successfully"}, nil
}
