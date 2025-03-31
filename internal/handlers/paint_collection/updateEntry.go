package paint_collection

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"paint-api/internal/db"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

type updateCollectionEntryInputBody struct {
	Quantity int `json:"quantity"`
	PaintId  int `json:"paint_id"`
}

type updateCollectionEntryInput struct {
	Id   int `path:"id"`
	Body updateCollectionEntryInputBody
}

type updateCollectionEntryOutputBody struct {
	Entry db.PaintCollection `json:"paint"`
}

type updateCollectionEntryOutput struct {
	Body db.CollectionPaintDetails `json:"body"`
}

var UpdateCollectionEntryOperation = huma.Operation{
	Method: http.MethodPut,
	Path:   "/collection/{id}",
	Tags:   []string{"collection"},
}

func UpdateCollectionEntryHandler(ctx context.Context, input *updateCollectionEntryInput) (*updateCollectionEntryOutput, error) {
	out := updateCollectionEntryOutput{
		Body: db.CollectionPaintDetails{},
	}
	connection, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return nil, errors.New("could not retrieve db from context")
	}

	err := verifyCollectionOwnership(ctx, connection, input.Id)
	if err != nil {
		return nil, err
	}

	var entry db.PaintCollection
	err = connection.First(&entry, input.Id).Error
	if err != nil {
		return nil, err
	}

	entry.Quantity = input.Body.Quantity
	entry.PaintId = input.Body.PaintId
	entry.UpdatedAt = time.Now()
	connection.Save(&entry)

	err = connection.Preload("Paint").
		Preload("Paint.Brand").
		First(&out.Body, entry.ID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, huma.NewError(http.StatusNotFound, "user not found")
		}
		slog.Error("An error occurred when fetching collection.", "error", err)
		return nil, huma.NewError(http.StatusInternalServerError, "could not fetch collection")
	}

	return &out, nil
}
