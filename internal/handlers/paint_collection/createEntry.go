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

type addToCollectionInputBody struct {
	PaintId  int `json:"paint_id" validate:"required"`
	Quantity int `json:"quantity" validate:"required"`
}
type addToCollectionInput struct {
	Body addToCollectionInputBody
}
type addToCollectionOutput struct {
	Body db.CollectionPaintDetails `json:"body"`
}

var AddToCollectionOperation = huma.Operation{
	Method: http.MethodPost,
	Path:   "/collection",
	Tags:   []string{"collection"},
}

func AddToCollectionHandler(ctx context.Context, input *addToCollectionInput) (*addToCollectionOutput, error) {
	out := addToCollectionOutput{
		Body: db.CollectionPaintDetails{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Quantity:  input.Body.Quantity,
			PaintID:   input.Body.PaintId,
		},
	}
	connection, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return nil, errors.New("could not retrieve db from context")
	}
	userId, ok := ctx.Value("userId").(string)
	if !ok {
		return nil, errors.New("could not retrieve user_id from context")
	}

	user := db.Users{}
	err := connection.Where("google_user_id = ?", userId).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, huma.NewError(http.StatusNotFound, "user not found")
		}
		slog.Error("An error occurred when fetching user.", "error", err)
		return nil, huma.NewError(http.StatusInternalServerError, "could not add paint to collection")
	}

	collectionEntry := db.PaintCollection{
		UserId:    user.ID,
		PaintId:   out.Body.PaintID,
		Quantity:  out.Body.Quantity,
		CreatedAt: out.Body.CreatedAt,
		UpdatedAt: out.Body.UpdatedAt,
	}

	err = connection.Create(&collectionEntry).Error
	if err != nil {
		slog.Error("An error occurred when creating entry.", "error", err)
		return nil, huma.NewError(http.StatusInternalServerError, "could not add paint to collection")
	}
	slog.Info("Paint added to collection", "Id", collectionEntry.ID, "PaintId", collectionEntry.PaintId, "Quantity", collectionEntry.Quantity)

	err = connection.Preload("Paint").Preload("Paint.Brand").First(&out.Body, collectionEntry.ID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, huma.NewError(http.StatusNotFound, "user not found")
		}
		slog.Error("An error occurred when fetching collection.", "error", err)
		return nil, huma.NewError(http.StatusInternalServerError, "could not fetch collection")
	}

	return &out, nil
}
