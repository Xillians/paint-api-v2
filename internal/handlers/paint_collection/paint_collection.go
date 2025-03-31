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

func RegisterRoutes(api huma.API) {
	huma.Register(api, createOperation, createHandler)
	huma.Register(api, listOperation, listHandler)
	huma.Register(api, updateOperation, updateHandler)
	huma.Register(api, deleteOperation, deleteHandler)
}

// verifyCollectionOwnership checks if the user owns the entry in the collection
//
// returns an error if the user does not own the entry.
func verifyCollectionOwnership(ctx context.Context, connection *gorm.DB, collectionId int) error {
	userId, ok := ctx.Value("userId").(string)
	if !ok {
		return huma.NewError(http.StatusNotFound, "Entry not found")
	}

	_, err := db.CollectionPaintDetails{}.GetEntry(connection, collectionId, userId)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return huma.NewError(http.StatusNotFound, "Entry not found")
		}
		slog.Error("Error getting entry", "error", err)
		return err
	}
	return nil
}
