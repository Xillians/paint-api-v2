package brands

import (
	"context"
	"errors"
	"net/http"
	"paint-api/internal/db"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

type getBrandInput struct {
	ID uint `path:"id" example:"1" required:"true"`
}

type getBrandOutput struct {
	Body db.PaintBrands
}

var GetOperation = huma.Operation{
	Method: http.MethodGet,
	Path:   "/paint-brands/{id}",
	Tags:   []string{"paint-brands"},
}

func GetHandler(ctx context.Context, input *getBrandInput) (*getBrandOutput, error) {
	connection, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return nil, errors.New("could not retrieve db from context")
	}

	brand, err := db.PaintBrands{}.GetBrand(connection, int(input.ID))
	if err != nil {
		return nil, err
	}

	return &getBrandOutput{Body: *brand}, nil
}
