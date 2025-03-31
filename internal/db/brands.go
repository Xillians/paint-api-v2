package db

import (
	"log/slog"
	"time"

	"gorm.io/gorm"
)

type PaintBrands struct {
	ID        int    `json:"id" gorm:"primaryKey"`
	Name      string `json:"name" gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (b PaintBrands) TableName() string {
	return "paint_brands"
}

type CreateBrandInput struct {
	Name string `json:"name"`
}
type UpdateBrandInput struct {
	Name string `json:"name"`
}

func (b PaintBrands) CreateBrand(connection *gorm.DB, input *CreateBrandInput) (*PaintBrands, error) {
	brand := PaintBrands{Name: input.Name}
	tx := connection.Create(&brand)
	if tx.Error != nil {
		slog.Error("Failed to create brand", "error", tx.Error, "transaction", tx, "brand", brand)
		return nil, tx.Error
	}
	return &brand, nil
}

func (b PaintBrands) ListBrands(connection *gorm.DB) ([]PaintBrands, error) {
	var brands []PaintBrands
	tx := connection.Find(&brands)
	if tx.Error != nil {
		slog.Error("Failed to list brands", "error", tx.Error, "transaction", tx)
		return nil, tx.Error
	}
	return brands, nil
}

func (b PaintBrands) GetBrand(connection *gorm.DB, id int) (*PaintBrands, error) {
	var brand PaintBrands
	tx := connection.First(&brand, id)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			slog.Info("Brand not found", "id", id)
			return nil, ErrRecordNotFound
		}
		slog.Error("Failed to get brand", "error", tx.Error, "transaction", tx)
		return nil, tx.Error
	}
	return &brand, nil
}

func (b PaintBrands) UpdateBrand(connection *gorm.DB, id int, input *UpdateBrandInput) (*PaintBrands, error) {
	var brand PaintBrands
	tx := connection.First(&brand, id)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			slog.Info("Brand not found", "id", id)
			return nil, ErrRecordNotFound
		}
		slog.Error("Failed to get brand", "error", tx.Error, "transaction", tx)
		return nil, tx.Error
	}

	brand.Name = input.Name
	tx = connection.Save(&brand)
	if tx.Error != nil {
		slog.Error("Failed to update brand", "error", tx.Error, "transaction", tx)
		return nil, tx.Error
	}
	return &brand, nil
}

func (b PaintBrands) DeleteBrand(connection *gorm.DB, id int) error {
	var brand PaintBrands
	tx := connection.First(&brand, id)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			slog.Info("Brand not found", "id", id)
			return ErrRecordNotFound
		}
		slog.Error("Failed to get brand", "error", tx.Error, "transaction", tx)
		return tx.Error
	}
	tx = connection.Delete(&brand)
	if tx.Error != nil {
		slog.Error("Failed to delete brand", "error", tx.Error, "transaction", tx)
		return tx.Error
	}
	return nil
}
