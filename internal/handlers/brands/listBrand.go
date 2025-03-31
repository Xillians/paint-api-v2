package brands

import (
	"context"
	"errors"
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

var ListOperation = huma.Operation{
	Method: http.MethodGet,
	Path:   "/paint-brands",
	Tags:   []string{"paint-brands"},
}

func ListHandler(ctx context.Context, input *listBrandInput) (*listBrandOutput, error) {
	connection, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return nil, errors.New("could not retrieve db from context")
	}

	brands, err := db.PaintBrands{}.ListBrands(connection)
	if err != nil {
		return nil, err
	}

	return &listBrandOutput{Body: listBrandOutputBody{Brands: brands}}, nil
}
