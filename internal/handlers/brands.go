package handlers

import (
	"context"
	"errors"
	"net/http"
	"paint-api/internal/db"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

type CreateBrandInputBody struct {
	Name string `json:"name"`
}
type CreatebrandInput struct {
	Body CreateBrandInputBody
}

type CreateBrandOutput struct {
	Body db.PaintBrands
}

var CreatePaintBrandOperation = huma.Operation{
	Method: http.MethodPost,
	Path:   "/paint-brands",
}

func CreatePaintBrandHandler(ctx context.Context, input *CreatebrandInput) (*CreateBrandOutput, error) {
	connection, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return nil, errors.New("could not retrieve db from context")
	}
	brand := db.PaintBrands{Name: input.Body.Name}
	connection.Create(&brand)
	return &CreateBrandOutput{Body: brand}, nil
}

type GetBrandInput struct {
}
type GetBrandOutput struct {
	Body []db.PaintBrands
}

var GetPaintBrandsOperation = huma.Operation{
	Method: http.MethodGet,
	Path:   "/paint-brands",
}

func GetPaintBrandsHandler(ctx context.Context, input *GetBrandInput) (*GetBrandOutput, error) {
	connection, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return nil, errors.New("could not retrieve db from context")
	}
	var brands []db.PaintBrands
	connection.Find(&brands)
	return &GetBrandOutput{Body: brands}, nil
}

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

var UpdatePaintBrandOperation = huma.Operation{
	Method: http.MethodPut,
	Path:   "/paint-brands/{id}",
}

func UpdatePaintBrandHandler(ctx context.Context, input *UpdateBrandInput) (*UpdateBrandOutput, error) {
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
	return &UpdateBrandOutput{Body: brand}, nil
}

type DeleteBrandInput struct {
	ID uint `path:"id" example:"1" required:"true"`
}

type DeleteBrandOutput struct {
	Body string
}

var DeletePaintBrandOperation = huma.Operation{
	Method: http.MethodDelete,
	Path:   "/paint-brands/{id}",
}

func DeletePaintBrandHandler(ctx context.Context, input *DeleteBrandInput) (*DeleteBrandOutput, error) {
	connection, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return nil, errors.New("could not retrieve db from context")
	}
	if err := connection.Delete(&db.PaintBrands{}, input.ID).Error; err != nil {
		return nil, err
	}
	return &DeleteBrandOutput{Body: "Brand deleted successfully"}, nil
}
