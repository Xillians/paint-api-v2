package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"paint-api/internal/db"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

type PaintOutputDetails struct {
	PaintId        int            `json:"id" gorm:"primaryKey"`
	PaintName      string         `json:"name"`
	BrandId        int            `json:"-" gorm:"not null"`
	ColorCode      string         `json:"color_code"`
	Description    string         `json:"description"`
	PaintCreatedAt time.Time      `json:"created_at"`
	PaintUpdatedAt time.Time      `json:"updated_at"`
	Brand          db.PaintBrands `json:"brand" gorm:"foreignKey:BrandId"`
}

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
	Body PaintOutputDetails
}

var CreatePaintOperation = huma.Operation{
	Method: http.MethodPost,
	Path:   "/paints",
	Tags:   []string{"paints"},
}

func CreatePaintHandler(ctx context.Context, input *createPaintInput) (*createPaintOutput, error) {
	out := createPaintOutput{
		Body: PaintOutputDetails{
			PaintName:      input.Body.Name,
			ColorCode:      input.Body.ColorCode,
			Description:    input.Body.Description,
			BrandId:        input.Body.BrandId,
			PaintCreatedAt: time.Now(),
			PaintUpdatedAt: time.Now(),
		},
	}
	connection, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return nil, errors.New("could not retrieve db from context")
	}
	paint := db.Paints{
		Name:        out.Body.PaintName,
		BrandId:     out.Body.BrandId,
		ColorCode:   out.Body.ColorCode,
		Description: out.Body.Description,
	}
	connection.Create(&paint)

	err := connection.Preload("Brand").
		Table("paints").
		Select(`
			paints.id AS paint_id, 
			paints.name as paint_name, 
			paints.created_at as paint_created_at, 
			paints.updated_at as paint_updated_at
		`).
		Where("paints.id = ?", paint.ID).
		Find(&out.Body).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, huma.NewError(http.StatusNotFound, "paint not found")
		}
		slog.Error("An error occurred when fetching paint.", "error", err)
		return nil, huma.NewError(http.StatusInternalServerError, "could not fetch paint")
	}
	return &out, nil
}

type listPaintInput struct {
}

type listPaintOutputBody struct {
	Paints []PaintOutputDetails `json:"paints"`
}

type listPaintOutput struct {
	Body listPaintOutputBody `json:"body"`
}

var ListPaintsOperation = huma.Operation{
	Method: http.MethodGet,
	Path:   "/paints",
	Tags:   []string{"paints"},
}

func ListPaintsHandler(ctx context.Context, input *listPaintInput) (*listPaintOutput, error) {
	out := listPaintOutput{
		Body: listPaintOutputBody{
			Paints: []PaintOutputDetails{},
		},
	}
	connection, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return nil, errors.New("could not retrieve db from context")
	}

	err := connection.Preload("Brand").
		Table("paints").
		Select(`
			paints.id AS paint_id,
			paints.name as paint_name,
			paints.created_at as paint_created_at,
			paints.updated_at as paint_updated_at,
			paints.description,
			paints.color_code,
			paints.brand_id
		`).
		Joins("JOIN paint_brands ON paint_brands.id = paints.brand_id").
		Find(&out.Body.Paints).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, huma.NewError(http.StatusNotFound, "paint not found")
		}
		slog.Error("An error occurred when fetching paint.", "error", err)
		return nil, huma.NewError(http.StatusInternalServerError, "could not fetch paint")
	}

	return &out, nil
}

type getPaintsInput struct {
	Id int `path:"id"`
}

type getPaintOutput struct {
	Body PaintOutputDetails `json:"body"`
}

var GetPaintsOperation = huma.Operation{
	Method: http.MethodGet,
	Path:   "/paints/{id}",
	Tags:   []string{"paints"},
}

func GetPaintHandler(ctx context.Context, input *getPaintsInput) (*getPaintOutput, error) {
	out := getPaintOutput{
		Body: PaintOutputDetails{},
	}
	connection, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return nil, errors.New("could not retrieve db from context")
	}

	err := connection.Preload("Brand").
		Table("paints").
		Select(`
			paints.id AS paint_id,
			paints.name as paint_name,
			paints.created_at as paint_created_at,
			paints.updated_at as paint_updated_at,
			paints.description,
			paints.color_code,
			paints.brand_id
		`).
		Where("paints.id = ?", input.Id).
		Find(&out.Body).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, huma.NewError(http.StatusNotFound, "paint not found")
		}
		slog.Error("An error occurred when fetching paint.", "error", err)
		return nil, huma.NewError(http.StatusInternalServerError, "could not fetch paint")
	}

	return &out, nil
}

type updatePaintInputBody struct {
	Name string `json:"name"`
}
type updatePaintInput struct {
	Id   int `path:"id"`
	Body updatePaintInputBody
}

type updatePaintOutput struct {
	Body PaintOutputDetails
}

var UpdatePaintOperation = huma.Operation{
	Method: http.MethodPut,
	Path:   "/paints/{id}",
	Tags:   []string{"paints"},
}

func UpdatePaintHandler(ctx context.Context, input *updatePaintInput) (*updatePaintOutput, error) {
	out := updatePaintOutput{
		Body: PaintOutputDetails{},
	}

	connection, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return nil, errors.New("could not retrieve db from context")
	}
	var paint db.Paints
	if err := connection.First(&paint, input.Id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, huma.NewError(http.StatusNotFound, "paint not found")
		}
		slog.Error("An error occurred when fetching paint.", "error", err)
		return nil, huma.NewError(http.StatusInternalServerError, "could not fetch paint")
	}

	paint.Name = input.Body.Name
	connection.Save(&paint)

	err := connection.Preload("Brand").
		Table("paints").
		Select(`
			paints.id AS paint_id,
			paints.name as paint_name,
			paints.created_at as paint_created_at,
			paints.updated_at as paint_updated_at,
			paints.description,
			paints.color_code,
			paints.brand_id
		`).
		Joins("JOIN paint_brands ON paint_brands.id = paints.brand_id").
		Find(&out.Body).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, huma.NewError(http.StatusNotFound, "paint not found")
		}
		slog.Error("An error occurred when fetching paint.", "error", err)
		return nil, huma.NewError(http.StatusInternalServerError, "could not fetch paint")
	}

	return &out, nil
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
