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

type AddToCollectionInputBody struct {
	Quantity int `json:"quantity"`
	PaintID  int `json:"paint_id"`
}

type AddToCollectionInput struct {
	Body AddToCollectionInputBody `json:"body"`
}
type AddToCollectionOutput struct {
	Body db.CollectionPaintDetails `json:"body"`
}

var createOperation = huma.Operation{
	Method: http.MethodPost,
	Path:   "/collection",
	Tags:   []string{"collection"},
}

func CreateHandler(ctx context.Context, input *AddToCollectionInput) (*AddToCollectionOutput, error) {
	connection, ok := ctx.Value(middleware.DbKey).(*gorm.DB)
	if !ok {
		slog.Error("could not retrieve db from context")
		return nil, huma.NewError(http.StatusInternalServerError, "could not add paint to collection")
	}
	userId, ok := ctx.Value(middleware.UserIdKey).(string)
	if !ok {
		slog.Error("could not retrieve userId from context")
		return nil, huma.NewError(http.StatusInternalServerError, "could not add paint to collection")
	}

	user, err := db.Users{}.GetUserByGoogleId(connection, userId)
	if err != nil {
		slog.Error("failed to get user by google id", "error", err)
		return nil, huma.NewError(http.StatusInternalServerError, "could not add paint to collection")
	}

	entry, err := db.CollectionPaintDetails{}.CreateEntry(
		connection, db.CreateCollectionEntryInput{
			UserId:   user.ID,
			PaintID:  input.Body.PaintID,
			Quantity: input.Body.Quantity,
		})
	if err != nil {
		slog.Error("failed to create collection entry", "error", err)
		return nil, huma.NewError(http.StatusInternalServerError, "could not add paint to collection")
	}

	return &AddToCollectionOutput{Body: *entry}, nil
}
