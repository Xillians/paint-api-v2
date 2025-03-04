package handlers

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

var CreatePaintBrandOperation = huma.Operation{
	Method: http.MethodPost,
	Path:   "/paint-brands",
	Tags:   []string{"paint-brands"},
}

func CreatePaintBrandHandler(ctx context.Context, input *createbrandInput) (*createBrandOutput, error) {
	connection, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return nil, errors.New("could not retrieve db from context")
	}
	brand := db.PaintBrands{Name: input.Body.Name}
	connection.Create(&brand)
	return &createBrandOutput{Body: brand}, nil
}

type listBrandInput struct {
}

type listBrandOutputBody struct {
	Brands []db.PaintBrands `json:"brands"`
}

type listBrandOutput struct {
	Body listBrandOutputBody `json:"body"`
}

var ListPaintBrandsOperation = huma.Operation{
	Method: http.MethodGet,
	Path:   "/paint-brands",
	Tags:   []string{"paint-brands"},
}

func ListPaintBrandsHandler(ctx context.Context, input *listBrandInput) (*listBrandOutput, error) {
	out := listBrandOutput{
		Body: listBrandOutputBody{
			Brands: []db.PaintBrands{},
		},
	}
	connection, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return nil, errors.New("could not retrieve db from context")
	}
	connection.Find(&out.Body.Brands)
	return &out, nil
}

type getBrandInput struct {
	ID uint `path:"id" example:"1" required:"true"`
}

type getBrandOutput struct {
	Body db.PaintBrands
}

var GetPaintBrandOperation = huma.Operation{
	Method: http.MethodGet,
	Path:   "/paint-brands/{id}",
	Tags:   []string{"paint-brands"},
}

func GetPaintBrandHandler(ctx context.Context, input *getBrandInput) (*getBrandOutput, error) {
	connection, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return nil, errors.New("could not retrieve db from context")
	}
	var brand db.PaintBrands
	if err := connection.First(&brand, input.ID).Error; err != nil {
		return nil, err
	}
	return &getBrandOutput{Body: brand}, nil
}

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

var UpdatePaintBrandOperation = huma.Operation{
	Method: http.MethodPut,
	Path:   "/paint-brands/{id}",
	Tags:   []string{"paint-brands"},
}

func UpdatePaintBrandHandler(ctx context.Context, input *updateBrandInput) (*updateBrandOutput, error) {
	connection, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return nil, errors.New("could not retrieve db from context")
	}
	var brand db.PaintBrands
	if err := connection.First(&brand, input.ID).Error; err != nil {
		return nil, err
	}
	brand.Name = input.Body.Name
	connection.Save(&brand)
	return &updateBrandOutput{Body: brand}, nil
}

type deleteBrandInput struct {
	ID uint `path:"id" example:"1" required:"true"`
}

type deleteBrandOutput struct {
	Body string
}

var DeletePaintBrandOperation = huma.Operation{
	Method: http.MethodDelete,
	Path:   "/paint-brands/{id}",
	Tags:   []string{"paint-brands"},
}

func DeletePaintBrandHandler(ctx context.Context, input *deleteBrandInput) (*deleteBrandOutput, error) {
	connection, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return nil, errors.New("could not retrieve db from context")
	}
	if err := connection.Delete(&db.PaintBrands{}, input.ID).Error; err != nil {
		return nil, err
	}
	return &deleteBrandOutput{Body: "Brand deleted successfully"}, nil
}
