package brands

import (
	"context"
	"errors"
	"net/http"
	"paint-api/internal/db"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

type createBrandInputBody struct {
	Name string `json:"name"`
}
type createbrandInput struct {
	Body createBrandInputBody
}

type createBrandOutput struct {
	Body db.PaintBrands
}

var CreateOperation = huma.Operation{
	Method: http.MethodPost,
	Path:   "/paint-brands",
	Tags:   []string{"paint-brands"},
}

func CreateHandler(ctx context.Context, input *createbrandInput) (*createBrandOutput, error) {
	connection, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return nil, errors.New("could not retrieve db from context")
	}

	brand, err := db.PaintBrands{}.CreateBrand(connection, &db.CreateBrandInput{Name: input.Body.Name})
	if err != nil {
		huma.NewError(500, "Failed to create brand")
	}

	return &createBrandOutput{Body: *brand}, nil
}
