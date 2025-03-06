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
	slog.Debug("Creating entry", "entry", collectionEntry)
	if err := connection.Create(&collectionEntry).Error; err != nil {
		slog.Error("An error occurred when creating entry.", "error", err)
		return nil, huma.NewError(http.StatusInternalServerError, "could not add paint to collection")
	}
	return &addToCollectionOutput{Body: collectionEntry}, nil
}

type listPaintCollectionInput struct {
}

type listPaintCollectionOutputBody struct {
	Collection []db.PaintCollection `json:"collection"`
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
			Collection: []db.PaintCollection{},
		},
	}
	connection, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return nil, errors.New("could not retrieve db from context")
	}
	connection.Find(&out.Body.Collection)
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
	if err := connection.First(&out.Body.Paint, input.Id).Error; err != nil {
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
	if err := connection.First(&out.Body.Entry, input.Id).Error; err != nil {
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

func DeleteCollectionEntryHandler(ctx context.Context, input *deleteCollectionEntryInput) (*deleteCollectionEntryOutput, error) {
	connection, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return nil, errors.New("could not retrieve db from context")
	}
	if err := connection.Delete(&db.PaintCollection{}, input.Id).Error; err != nil {
		return nil, err
	}
	return &deleteCollectionEntryOutput{Body: "Paint deleted successfully"}, nil
}
