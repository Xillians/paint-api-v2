package db

import "time"

type Paints struct {
	ID          int       `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name"`
	BrandId     int       `json:"brand_id" foreignKey:"Brands(ID)"`
	ColorCode   string    `json:"color_code"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type PaintOutputDetails struct {
	Id          int         `json:"id" gorm:"primaryKey"`
	Name        string      `json:"name"`
	BrandId     int         `json:"-" gorm:"not null"`
	ColorCode   string      `json:"color_code"`
	Description string      `json:"description"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	Brand       PaintBrands `json:"brand" gorm:"foreignKey:BrandId"`
}

func (p PaintOutputDetails) TableName() string {
	return "paints"
}
