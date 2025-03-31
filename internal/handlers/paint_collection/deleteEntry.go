package paint_collection

import (
	"context"
	"errors"
	"net/http"
	"paint-api/internal/db"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

type deleteCollectionEntryInput struct {
	Id int `path:"id"`
}

type deleteCollectionEntryOutput struct {
	Body string
}

var DeleteCollectionEntryOperation = huma.Operation{
	Method: http.MethodDelete,
	Path:   "/collection/{id}",
	Tags:   []string{"collection"},
}

// DeleteCollectionEntryHandler deletes a paint from the collection
//
// Deletes a paint from the collection. Requires authentication, and the paint must belong to the user.
func DeleteCollectionEntryHandler(ctx context.Context, input *deleteCollectionEntryInput) (*deleteCollectionEntryOutput, error) {
	connection, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return nil, errors.New("could not retrieve db from context")
	}

	err := verifyCollectionOwnership(ctx, connection, input.Id)
	if err != nil {
		return nil, err
	}

	err = connection.Delete(&db.PaintCollection{}, input.Id).Error
	if err != nil {
		return nil, err
	}

	return &deleteCollectionEntryOutput{Body: "Paint deleted successfully"}, nil
}
