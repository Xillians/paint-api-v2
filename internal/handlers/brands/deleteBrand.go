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

type DeleteBrandInput struct {
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

func DeleteHandler(ctx context.Context, input *DeleteBrandInput) (*DeleteBrandOutput, error) {
	userRole := ctx.Value(middleware.RoleKey).(string)
	if userRole != "administrator" {
		return nil, huma.NewError(http.StatusForbidden, "You are not allowed to perform this action")
	}
	connection, ok := ctx.Value(middleware.DbKey).(*gorm.DB)
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
