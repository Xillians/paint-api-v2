package db

import (
	"log/slog"

	"gorm.io/gorm"
)

type Users struct {
	ID           int    `json:"id" gorm:"primaryKey"`
	GoogleUserId string `json:"google_user_id"`
	Email        string `json:"email"`
	Role         string `json:"role"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

func (u Users) TableName() string {
	return "users"
}

type RegisterUserInput struct {
	GoogleUserId string `json:"user_id" validate:"required"`
	Email        string `json:"email" required:"false"`
}

func (u Users) RegisterUser(connection *gorm.DB, input RegisterUserInput, role string) (*Users, error) {
	user := connection.First(&Users{}, "google_user_id = ?", input.GoogleUserId)
	if user.RowsAffected > 0 {
		return nil, ErrRecordExists
	}

	User := Users{
		GoogleUserId: input.GoogleUserId,
		Email:        input.Email,
		Role:         role,
		CreatedAt:    "time.Now().String()",
		UpdatedAt:    "time.Now().String()",
	}
	tx := connection.Create(&User)
	if tx.Error != nil {
		slog.Error("Failed to create user.", "error", tx.Error)
		return nil, tx.Error
	}

	return &User, nil
}

func (u Users) GetUserByGoogleId(connection *gorm.DB, googleUserId string) (*Users, error) {
	user := &Users{}
	tx := connection.First(&user, "google_user_id = ?", googleUserId)
	if tx.Error != nil {
		if tx.Error == gorm.ErrRecordNotFound {
			return nil, ErrRecordNotFound
		}
		return nil, tx.Error
	}

	return user, nil
}

func (u Users) DeleteUserByGoogleId(connection *gorm.DB, googleUserId string) error {
	user := &Users{}
	result := connection.First(&user, "google_user_id = ?", googleUserId)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return ErrRecordNotFound
		}
		slog.Error("Failed to find user.", "error", result.Error)
		return result.Error
	}

	tx := connection.Delete(&user)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}
