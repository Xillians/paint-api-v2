package brands

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"paint-api/internal/db"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

type deleteBrandInput struct {
	ID uint `path:"id" example:"1" required:"true"`
}

type DeleteBrandOutput struct {
	Body string
}

var deleteOperation = huma.Operation{
	Method: http.MethodDelete,
	Path:   "/paint-brands/{id}",
	Tags:   []string{"paint-brands"},
}

func deleteHandler(ctx context.Context, input *deleteBrandInput) (*DeleteBrandOutput, error) {
	userRole := ctx.Value("role").(string)
	if userRole != "administrator" {
		return nil, huma.NewError(http.StatusForbidden, "You are not allowed to perform this action")
	}
	connection, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		slog.Error("Could not retrieve db from context")
		return nil, huma.NewError(http.StatusInternalServerError, "failed to delete brand")
	}

	err := db.PaintBrands{}.DeleteBrand(connection, int(input.ID))
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, huma.NewError(http.StatusNotFound, "Brand not found")
		}
		slog.Error("Failed to delete brand", "error", err, "id", input.ID)
		return nil, huma.NewError(http.StatusInternalServerError, "Failed to delete brand")
	}

	return &DeleteBrandOutput{Body: "Brand deleted successfully"}, nil
}
