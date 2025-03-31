package db

import "time"

type PaintBrands struct {
	ID        int    `json:"id" gorm:"primaryKey"`
	Name      string `json:"name"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (b PaintBrands) TableName() string {
	return "brands"
}
