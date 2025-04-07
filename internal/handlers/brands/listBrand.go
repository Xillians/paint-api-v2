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

type ListBrandInput struct {
}

type ListBrandOutputBody struct {
	Brands []db.PaintBrands `json:"brands"`
}

type ListBrandOutput struct {
	Body ListBrandOutputBody `json:"body"`
}

var listOperation = huma.Operation{
	Method: http.MethodGet,
	Path:   "/paint-brands",
	Tags:   []string{"paint-brands"},
}

func ListHandler(ctx context.Context, input *ListBrandInput) (*ListBrandOutput, error) {
	connection, ok := ctx.Value(middleware.DbKey).(*gorm.DB)
	if !ok {
		slog.Error("Could not retrieve db from context")
		return nil, huma.NewError(http.StatusInternalServerError, "failed to list brands")
	}

	brands, err := db.PaintBrands{}.ListBrands(connection)
	if err != nil {
		slog.Error("Failed to list brands", "error", err)
		return nil, huma.NewError(http.StatusInternalServerError, "failed to list brands")
	}

	return &ListBrandOutput{Body: ListBrandOutputBody{Brands: brands}}, nil
}
