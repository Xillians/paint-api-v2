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

type UpdateBrandInputBody struct {
	Name string `json:"name"`
}
type UpdateBrandInput struct {
	ID   uint `path:"id" example:"1" required:"true"`
	Body UpdateBrandInputBody
}

type UpdateBrandOutput struct {
	Body db.PaintBrands
}

var updateOperation = huma.Operation{
	Method: http.MethodPut,
	Path:   "/paint-brands/{id}",
	Tags:   []string{"paint-brands"},
}

func UpdateHandler(ctx context.Context, input *UpdateBrandInput) (*UpdateBrandOutput, error) {
	userRole := ctx.Value(middleware.RoleKey).(string)
	if userRole != "administrator" {
		return nil, huma.NewError(http.StatusForbidden, "You are not allowed to perform this action")
	}

	connection, ok := ctx.Value(middleware.DbKey).(*gorm.DB)
	if !ok {
		slog.Error("Could not retrieve db from context")
		return nil, huma.NewError(http.StatusInternalServerError, "failed to update brand")
	}

	brand, err := db.PaintBrands{}.UpdateBrand(
		connection,
		int(input.ID),
		&db.UpdateBrandInput{Name: input.Body.Name},
	)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, huma.NewError(http.StatusNotFound, "Brand not found")
		}
		slog.Error("Failed to update brand", "error", err, "id", input.ID)
		return nil, huma.NewError(http.StatusInternalServerError, "Failed to update brand")
	}
	return &UpdateBrandOutput{Body: *brand}, nil
}
