package handlers

import (
	"context"
	"errors"
	"net/http"

	"paint-api/internal/db"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

type createPaintInputBody struct {
	Name        string `json:"name" validate:"required"`
	BrandId     int    `json:"brand_id" validate:"required"`
	ColorCode   string `json:"color_code" validate:"required"`
	Description string `json:"description"`
}
type createPaintInput struct {
	Body createPaintInputBody
}

type createPaintOutput struct {
	Body db.Paints
}

var CreatePaintOperation = huma.Operation{
	Method: http.MethodPost,
	Path:   "/paints",
	Tags:   []string{"paints"},
}

func CreatePaintHandler(ctx context.Context, input *createPaintInput) (*createPaintOutput, error) {
	connection, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return nil, errors.New("could not retrieve db from context")
	}
	paint := db.Paints{
		Name:        input.Body.Name,
		BrandId:     input.Body.BrandId,
		ColorCode:   input.Body.ColorCode,
		Description: input.Body.Description,
	}
	connection.Create(&paint)
	return &createPaintOutput{Body: paint}, nil
}

type listPaintInput struct {
}

type listPaintOutputBody struct {
	Paints []db.Paints `json:"paints"`
}

type listPaintOutput struct {
	Body listPaintOutputBody
}

var ListPaintsOperation = huma.Operation{
	Method: http.MethodGet,
	Path:   "/paints",
	Tags:   []string{"paints"},
}

func ListPaintsHandler(ctx context.Context, input *listPaintInput) (*listPaintOutput, error) {
	connection, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return nil, errors.New("could not retrieve db from context")
	}
	var paints []db.Paints
	connection.Find(&paints)
	return &listPaintOutput{Body: listPaintOutputBody{Paints: paints}}, nil
}

type getPaintsInput struct {
	Id int `path:"id"`
}

type getPaintOutput struct {
	Body db.Paints
}

var GetPaintsOperation = huma.Operation{
	Method: http.MethodGet,
	Path:   "/paints/{id}",
	Tags:   []string{"paints"},
}

func GetPaintHandler(ctx context.Context, input *getPaintsInput) (*getPaintOutput, error) {
	connection, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return nil, errors.New("could not retrieve db from context")
	}
	var paint db.Paints
	if err := connection.First(&paint, input.Id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, huma.NewError(http.StatusNotFound, "paint not found")
		}
		return nil, err
	}
	return &getPaintOutput{Body: paint}, nil
}

type updatePaintInputBody struct {
	Name string `json:"name"`
}
type updatePaintInput struct {
	Id   int `path:"id"`
	Body updatePaintInputBody
}

type updatePaintOutput struct {
	Body db.Paints
}

var UpdatePaintOperation = huma.Operation{
	Method: http.MethodPut,
	Path:   "/paints/{id}",
	Tags:   []string{"paints"},
}

func UpdatePaintHandler(ctx context.Context, input *updatePaintInput) (*updatePaintOutput, error) {
	connection, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return nil, errors.New("could not retrieve db from context")
	}
	var paint db.Paints
	if err := connection.First(&paint, input.Id).Error; err != nil {
		return nil, err
	}
	paint.Name = input.Body.Name
	connection.Save(&paint)
	return &updatePaintOutput{Body: paint}, nil
}

type deletePaintInput struct {
	Id int `path:"id"`
}

type deletePaintOutput struct {
	Body string
}

var DeletePaintOperation = huma.Operation{
	Method: http.MethodDelete,
	Path:   "/paints/{id}",
	Tags:   []string{"paints"},
}

func DeletePaintHandler(ctx context.Context, input *deletePaintInput) (*deletePaintOutput, error) {
	connection, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return nil, errors.New("could not retrieve db from context")
	}
	if err := connection.Delete(&db.Paints{}, input.Id).Error; err != nil {
		return nil, err
	}
	return &deletePaintOutput{Body: "Paint deleted successfully"}, nil
}
