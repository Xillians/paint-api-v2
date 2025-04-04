package brands

import (
	"context"
	"log/slog"
	"net/http"
	"paint-api/internal/db"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

type createbrandInput struct {
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

func createHandler(ctx context.Context, input *createbrandInput) (*CreateBrandOutput, error) {
	connection, ok := ctx.Value("db").(*gorm.DB)
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
