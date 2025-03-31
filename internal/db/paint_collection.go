package db

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

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
	UserId    int       `json:"-" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Paint     Paints    `json:"paint" gorm:"foreignKey:PaintID"`
	User      Users     `json:"-" gorm:"foreignKey:UserId"`
}

func (c CollectionPaintDetails) TableName() string {
	return "paint_collections"
}

type CreateCollectionEntryInput struct {
	Quantity int `json:"quantity" validate:"required"`
	PaintID  int `json:"paint_id" validate:"required"`
	UserId   int `json:"user_id" validate:"required"`
}
type UpdateCollectionEntryInput struct {
	ID       int `json:"id" validate:"required"`
	Quantity int `json:"quantity"`
	PaintID  int `json:"paint_id"`
}

func (c CollectionPaintDetails) CreateEntry(connection *gorm.DB, input CreateCollectionEntryInput) (*CollectionPaintDetails, error) {
	entry := CollectionPaintDetails{
		Quantity:  input.Quantity,
		PaintID:   input.PaintID,
		UserId:    input.UserId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	tx := connection.Create(&entry)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &entry, nil
}

func (c CollectionPaintDetails) UpdateEntry(connection *gorm.DB, input UpdateCollectionEntryInput) (*CollectionPaintDetails, error) {
	entry := CollectionPaintDetails{}
	err := connection.First(&entry, input.ID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRecordNotFound
		}
		return nil, err
	}

	entry.Quantity = input.Quantity
	entry.PaintID = input.PaintID
	entry.UpdatedAt = time.Now()
	tx := connection.Save(&entry)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &entry, nil
}

// Lists all entries in the collection that belong to the user.
func (c CollectionPaintDetails) ListEntries(connection *gorm.DB, userId int) ([]CollectionPaintDetails, error) {
	var entries []CollectionPaintDetails
	tx := connection.Preload("Paint").Find(&entries).Where("user_id = ?", userId)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, ErrRecordNotFound
		}
		return nil, tx.Error
	}
	return entries, nil
}

func (c CollectionPaintDetails) DeleteEntry(connection *gorm.DB, id int) error {
	return nil
}
