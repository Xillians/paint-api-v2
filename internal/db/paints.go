package db

import (
	"log/slog"
	"time"

	"gorm.io/gorm"
)

type Paints struct {
	Id          int         `json:"id" gorm:"primaryKey"`
	Name        string      `json:"name"`
	BrandId     int         `json:"-" gorm:"not null"`
	ColorCode   string      `json:"color_code"`
	Description string      `json:"description"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	Brand       PaintBrands `json:"brand" gorm:"foreignKey:BrandId"`
}

func (p Paints) TableName() string {
	return "paints"
}

type CreatePaintInput struct {
	Name        string `json:"name" validate:"required"`
	BrandId     int    `json:"brand_id" validate:"required"`
	ColorCode   string `json:"color_code" validate:"required"`
	Description string `json:"description"`
}
type UpdatePaintInput struct {
	Name        string `json:"name"`
	BrandId     int    `json:"brand_id"`
	ColorCode   string `json:"color_code"`
	Description string `json:"description"`
}

func (p Paints) CreatePaint(connection *gorm.DB, input *CreatePaintInput) (*Paints, error) {
	paint := Paints{
		Name:        input.Name,
		BrandId:     input.BrandId,
		ColorCode:   input.ColorCode,
		Description: input.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	tx := connection.Create(&paint)
	if tx.Error != nil {
		slog.Error("Failed to create paint", "error", tx.Error, "transaction", tx, "paint", paint)
		return nil, tx.Error
	}
	return &paint, nil
}

func (p Paints) GetPaint(connection *gorm.DB, id int) (*Paints, error) {
	paint := Paints{}
	tx := connection.Preload("Brand").First(&paint, id)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			slog.Error("Paint not found", "error", tx.Error, "transaction", tx, "id", id)
			return nil, ErrRecordNotFound
		}
		slog.Error("Failed to fetch paint", "error", tx.Error, "transaction", tx, "id", id)
		return nil, tx.Error
	}
	return &paint, nil
}

func (p Paints) ListPaints(connection *gorm.DB) ([]Paints, error) {
	var paints []Paints
	tx := connection.Preload("Brand").Find(&paints)
	if tx.Error != nil {
		slog.Error("Failed to list paints", "error", tx.Error, "transaction", tx)
		return nil, tx.Error
	}
	return paints, nil
}

func (p Paints) UpdatePaint(connection *gorm.DB, id int, input *UpdatePaintInput) (*Paints, error) {
	paint := Paints{}
	tx := connection.First(&paint, id)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			slog.Error("Paint not found", "error", tx.Error, "transaction", tx, "id", id)
			return nil, ErrRecordNotFound
		}
		slog.Error("Failed to fetch paint", "error", tx.Error, "transaction", tx, "id", id)
		return nil, tx.Error
	}

	paint.Name = input.Name
	paint.BrandId = input.BrandId
	paint.ColorCode = input.ColorCode
	paint.Description = input.Description
	paint.UpdatedAt = time.Now()

	tx = connection.Save(&paint)
	if tx.Error != nil {
		slog.Error("Failed to update paint", "error", tx.Error, "transaction", tx, "paint", paint)
		return nil, tx.Error
	}

	return &paint, nil
}

func (p Paints) DeletePaint(connection *gorm.DB, id int) error {
	paint := Paints{}
	tx := connection.First(&paint, id)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			slog.Error("Paint not found", "error", tx.Error, "transaction", tx, "id", id)
			return ErrRecordNotFound
		}
		slog.Error("Failed to fetch paint", "error", tx.Error, "transaction", tx, "id", id)
		return tx.Error
	}

	tx = connection.Delete(&paint)
	if tx.Error != nil {
		slog.Error("Failed to delete paint", "error", tx.Error, "transaction", tx, "paint", paint)
		return tx.Error
	}

	return nil
}
