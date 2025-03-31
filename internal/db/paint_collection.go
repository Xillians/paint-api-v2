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

type CollectionPaintDetails struct {
	ID        int       `json:"id" gorm:"primaryKey"`
	Quantity  int       `json:"quantity"`
	PaintID   int       `json:"-" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Paint     Paints    `json:"paint" gorm:"foreignKey:PaintID"`
}

func (c CollectionPaintDetails) TableName() string {
	return "paint_collections"
}
