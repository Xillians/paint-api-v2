package handlers

import (
	"context"
	"errors"
	"net/http"
	"paint-api/internal/db"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

type PaintBrandInput struct {
	Name string `json:"name"`
}

type PaintBrandOutput struct {
	Body db.PaintBrands
}

var CreatePaintBrandOperation = huma.Operation{
	Method: http.MethodPost,
	Path:   "/paint-brands",
}

func CreatePaintBrandHandler(ctx context.Context, input *PaintBrandInput) (*PaintBrandOutput, error) {
	connection, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return nil, errors.New("could not retrieve db from context")
	}
	brand := db.PaintBrands{Name: input.Name}
	connection.Create(&brand)
	return &PaintBrandOutput{Body: brand}, nil
}

var GetPaintBrandsOperation = huma.Operation{
	Method: http.MethodGet,
	Path:   "/paint-brands",
}

func GetPaintBrandsHandler(ctx context.Context) ([]db.PaintBrands, error) {
	connection, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return nil, errors.New("could not retrieve db from context")
	}
	var brands []db.PaintBrands
	connection.Find(&brands)
	return brands, nil
}

var GetPaintBrandOperation = huma.Operation{
	Method: http.MethodGet,
	Path:   "/paint-brands/{id}",
}

func GetPaintBrandHandler(ctx context.Context, id string) (*PaintBrandOutput, error) {
	connection, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return nil, errors.New("could not retrieve db from context")
	}
	var brand db.PaintBrands
	connection.First(&brand, id)
	return &PaintBrandOutput{Body: brand}, nil
}

var UpdatePaintBrandOperation = huma.Operation{
	Method: http.MethodPut,
	Path:   "/paint-brands/{id}",
}

func UpdatePaintBrandHandler(ctx context.Context, id string, input *PaintBrandInput) (*PaintBrandOutput, error) {
	connection, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return nil, errors.New("could not retrieve db from context")
	}
	var brand db.PaintBrands
	connection.First(&brand, id)
	connection.Model(&brand).Updates(db.PaintBrands{Name: input.Name})
	return &PaintBrandOutput{Body: brand}, nil
}

var DeletePaintBrandOperation = huma.Operation{
	Method: http.MethodDelete,
	Path:   "/paint-brands/{id}",
}

func DeletePaintBrandHandler(ctx context.Context, id string) (struct{}, error) {
	connection, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return struct{}{}, errors.New("could not retrieve db from context")
	}
	connection.Delete(&db.PaintBrands{}, id)
	return struct{}{}, nil
}
