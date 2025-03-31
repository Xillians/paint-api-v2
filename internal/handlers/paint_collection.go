package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"paint-api/internal/db"

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
