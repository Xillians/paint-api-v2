package db

import "time"

type PaintCollection struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	UserId    int       `json:"user_id" foreignKey:"Users(ID)"`
	PaintId   int       `json:"paint_id" foreignKey:"Paints(ID)"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
