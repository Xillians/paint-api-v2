package paint_collection

import (
	"context"
	"errors"

	"net/http"

	"paint-api/internal/db"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

// verifyCollectionOwnership checks if the user owns the entry in the collection
//
// returns an error if the user does not own the entry.
func verifyCollectionOwnership(ctx context.Context, connection *gorm.DB, collectionId int) error {
	userId, ok := ctx.Value("userId").(string)
	if !ok {
		return errors.New("could not retrieve user_id from context")
	}

	entry := db.PaintCollection{}
	if err := connection.Joins("JOIN users ON users.id = paint_collections.user_id").
		Where("paint_collections.id = ? AND users.google_user_id = ?", collectionId, userId).
		First(&entry).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return huma.NewError(http.StatusNotFound, "entry not found")
		}
		return err
	}
	return nil
}
