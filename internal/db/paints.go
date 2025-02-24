package db

import "time"

type Paints struct {
	ID          int       `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name"`
	BrandId     int       `json:"brand_id"`
	ColorCode   string    `json:"color_code"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
