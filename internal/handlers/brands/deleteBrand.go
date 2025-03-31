package brands

import (
	"context"
	"errors"
	"net/http"
	"paint-api/internal/db"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

type deleteBrandInput struct {
	ID uint `path:"id" example:"1" required:"true"`
}

type deleteBrandOutput struct {
	Body string
}

var DeleteOperation = huma.Operation{
	Method: http.MethodDelete,
	Path:   "/paint-brands/{id}",
	Tags:   []string{"paint-brands"},
}

func DeleteHandler(ctx context.Context, input *deleteBrandInput) (*deleteBrandOutput, error) {
	userRole := ctx.Value("role").(string)
	if userRole != "administrator" {
		return nil, huma.NewError(403, "You are not allowed to perform this action")
	}
	connection, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return nil, errors.New("could not retrieve db from context")
	}

	err := db.PaintBrands{}.DeleteBrand(connection, int(input.ID))
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, huma.NewError(404, "Brand not found")
		}
		return nil, err
	}

	return &deleteBrandOutput{Body: "Brand deleted successfully"}, nil
}
