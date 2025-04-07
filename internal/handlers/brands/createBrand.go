package brands

import (
	"context"
	"log/slog"
	"net/http"
	"paint-api/internal/db"
	"paint-api/internal/middleware"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

type CreatebrandInput struct {
	Body db.CreateBrandInput `json:"body"`
}

type CreateBrandOutput struct {
	Body db.PaintBrands
}

var createOperation = huma.Operation{
	Method: http.MethodPost,
	Path:   "/paint-brands",
	Tags:   []string{"paint-brands"},
}

func CreateHandler(ctx context.Context, input *CreatebrandInput) (*CreateBrandOutput, error) {
	connection, ok := ctx.Value(middleware.DbKey).(*gorm.DB)
	if !ok {
		slog.Error("Could not retrieve db from context")
		return nil, huma.NewError(http.StatusInternalServerError, "failed to create brand")
	}

	brand, err := db.PaintBrands{}.CreateBrand(connection, &db.CreateBrandInput{Name: input.Body.Name})
	if err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, "Failed to create brand")
	}

	return &CreateBrandOutput{Body: *brand}, nil
}
