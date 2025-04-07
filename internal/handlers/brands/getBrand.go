package brands

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"paint-api/internal/db"
	"paint-api/internal/middleware"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

type GetBrandInput struct {
	ID uint `path:"id" example:"1" required:"true"`
}

type GetBrandOutput struct {
	Body db.PaintBrands
}

var getOperation = huma.Operation{
	Method: http.MethodGet,
	Path:   "/paint-brands/{id}",
	Tags:   []string{"paint-brands"},
}

func GetHandler(ctx context.Context, input *GetBrandInput) (*GetBrandOutput, error) {
	connection, ok := ctx.Value(middleware.DbKey).(*gorm.DB)
	if !ok {
		slog.Error("Could not retrieve db from context")
		return nil, huma.NewError(http.StatusInternalServerError, "failed to get brand")
	}

	brand, err := db.PaintBrands{}.GetBrand(connection, int(input.ID))
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, huma.NewError(http.StatusNotFound, "Brand not found")
		}
		slog.Error("Failed to get brand", "error", err, "id", input.ID)
		return nil, huma.NewError(http.StatusInternalServerError, "Failed to get brand")
	}

	return &GetBrandOutput{Body: *brand}, nil
}
