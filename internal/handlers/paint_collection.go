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
type addToCollectoinInput struct {
	Body addToCollectionInputBody
}

type addToCollectionOutput struct {
	Body db.PaintCollection
}

var AddToCollectionOperation = huma.Operation{
	Method: http.MethodPost,
	Path:   "/collection",
	Tags:   []string{"collection"},
}

func AddToCollectionHandler(ctx context.Context, input *addToCollectoinInput) (*addToCollectionOutput, error) {
	connection, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return nil, errors.New("could not retrieve db from context")
	}
	userId, ok := ctx.Value("userId").(string)
	if !ok {
		return nil, errors.New("could not retrieve user_id from context")
	}

	user := db.Users{}
	if err := connection.Where("google_user_id = ?", userId).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, huma.NewError(http.StatusNotFound, "user not found")
		}
		slog.Error("An error occurred when fetching user.", "error", err)
		return nil, huma.NewError(http.StatusInternalServerError, "could not add paint to collection")
	}

	collectionEntry := db.PaintCollection{
		UserId:    user.ID,
		PaintId:   input.Body.PaintId,
		Quantity:  input.Body.Quantity,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := connection.Create(&collectionEntry).Error; err != nil {
		slog.Error("An error occurred when creating entry.", "error", err)
		return nil, huma.NewError(http.StatusInternalServerError, "could not add paint to collection")
	}
	return &addToCollectionOutput{Body: collectionEntry}, nil
}

type listPaintCollectionInput struct {
}

type collectionPaintDetails struct {
	CollectionId        int       `json:"id" gorm:"primaryKey"`
	Quantity            int       `json:"quantity"`
	PaintID             int       `json:"-" gorm:"not null"`
	Paint               db.Paints `json:"paint" gorm:"foreignKey:PaintID;references:ID"`
	CollectionCreatedAt time.Time `json:"created_at"`
	CollectionUpdatedAt time.Time `json:"updated_at"`
}

type listPaintCollectionOutputBody struct {
	Collection []collectionPaintDetails `json:"collection"`
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
			Collection: []collectionPaintDetails{},
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

	err := connection.Preload("Paint").
		Table("paint_collections").
		Select("paint_collections.id AS collection_id, paint_collections.quantity, paint_collections.paint_id, paint_collections.created_at as collection_created_at, paint_collections.updated_at as collection_updated_at").
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

type getCollectionEntryInput struct {
	Id int `path:"id"`
}

type getCollectionEntryOutputBody struct {
	Paint db.PaintCollection `json:"paint"`
}

type getCollectionEntryOutput struct {
	Body getCollectionEntryOutputBody `json:"body"`
}

var GetCollectionEntryOperation = huma.Operation{
	Method: http.MethodGet,
	Path:   "/collection/{id}",
	Tags:   []string{"collection"},
}

func GetCollectionEntryHandler(ctx context.Context, input *getCollectionEntryInput) (*getCollectionEntryOutput, error) {
	out := getCollectionEntryOutput{
		Body: getCollectionEntryOutputBody{
			Paint: db.PaintCollection{},
		},
	}

	connection, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return nil, errors.New("could not retrieve db from context")
	}

	err := verifyCollectionOwnership(ctx, connection, input.Id)
	if err != nil {
		return nil, err
	}

	err = connection.First(&out.Body.Paint, input.Id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, huma.NewError(http.StatusNotFound, "entry not found")
		}
		return nil, err
	}

	return &out, nil
}

type updateCollectionEntryInputBody struct {
	Quantity int `json:"quantity"`
}

type updateCollectionEntryInput struct {
	Id   int `path:"id"`
	Body updateCollectionEntryInputBody
}

type updateCollectionEntryOutputBody struct {
	Entry db.PaintCollection `json:"paint"`
}

type updateCollectionEntryOutput struct {
	Body updateCollectionEntryOutputBody `json:"body"`
}

var UpdateCollectionEntryOperation = huma.Operation{
	Method: http.MethodPut,
	Path:   "/collection/{id}",
	Tags:   []string{"collection"},
}

func UpdateCollectionEntryHandler(ctx context.Context, input *updateCollectionEntryInput) (*updateCollectionEntryOutput, error) {
	out := updateCollectionEntryOutput{
		Body: updateCollectionEntryOutputBody{
			Entry: db.PaintCollection{},
		},
	}
	connection, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return nil, errors.New("could not retrieve db from context")
	}

	err := verifyCollectionOwnership(ctx, connection, input.Id)
	if err != nil {
		return nil, err
	}

	err = connection.First(&out.Body.Entry, input.Id).Error
	if err != nil {
		return nil, err
	}

	out.Body.Entry.Quantity = input.Body.Quantity
	out.Body.Entry.UpdatedAt = time.Now()
	connection.Save(&out.Body.Entry)

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
