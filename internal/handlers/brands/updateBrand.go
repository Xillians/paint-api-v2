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

type updateBrandInputBody struct {
	Name string `json:"name"`
}
type updateBrandInput struct {
	ID   uint `path:"id" example:"1" required:"true"`
	Body updateBrandInputBody
}

type updateBrandOutput struct {
	Body db.PaintBrands
}

var UpdateOperation = huma.Operation{
	Method: http.MethodPut,
	Path:   "/paint-brands/{id}",
	Tags:   []string{"paint-brands"},
}

func UpdateHandler(ctx context.Context, input *updateBrandInput) (*updateBrandOutput, error) {
	userRole := ctx.Value("role").(string)
	if userRole != "administrator" {
		return nil, huma.NewError(403, "You are not allowed to perform this action")
	}

	connection, ok := ctx.Value("db").(*gorm.DB)
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
			return nil, huma.NewError(404, "Brand not found")
		}
		slog.Error("Failed to update brand", "error", err, "id", input.ID)
		return nil, huma.NewError(http.StatusInternalServerError, "Failed to update brand")
	}
	return &updateBrandOutput{Body: *brand}, nil
}
