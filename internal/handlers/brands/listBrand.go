package brands

import (
	"context"
	"log/slog"
	"net/http"
	"paint-api/internal/db"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

type listBrandInput struct {
}

type listBrandOutputBody struct {
	Brands []db.PaintBrands `json:"brands"`
}

type listBrandOutput struct {
	Body listBrandOutputBody `json:"body"`
}

var listOperation = huma.Operation{
	Method: http.MethodGet,
	Path:   "/paint-brands",
	Tags:   []string{"paint-brands"},
}

func listHandler(ctx context.Context, input *listBrandInput) (*listBrandOutput, error) {
	connection, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		slog.Error("Could not retrieve db from context")
		return nil, huma.NewError(http.StatusInternalServerError, "failed to list brands")
	}

	brands, err := db.PaintBrands{}.ListBrands(connection)
	if err != nil {
		slog.Error("Failed to list brands", "error", err)
		return nil, huma.NewError(http.StatusInternalServerError, "failed to list brands")
	}

	return &listBrandOutput{Body: listBrandOutputBody{Brands: brands}}, nil
}
