package db

import (
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

func (u Users) RegisterUser(connection *gorm.DB, input RegisterUserInput) (*Users, error) {
	user := connection.First(&Users{}, "google_user_id = ?", input.GoogleUserId)
	if user.RowsAffected > 0 {
		return nil, ErrRecordExists
	}

	User := Users{
		GoogleUserId: input.GoogleUserId,
		Email:        input.Email,
		Role:         "user",
		CreatedAt:    "time.Now().String()",
		UpdatedAt:    "time.Now().String()",
	}
	connection.Create(&User)

	return &User, nil
}

func (u Users) GetUserByGoogleId(connection *gorm.DB, googleUserId string) (*Users, error) {
	user := &Users{}
	result := connection.First(&user, "google_user_id = ?", googleUserId)
	if result.RowsAffected == 0 {
		return nil, ErrRecordNotFound
	}

	return user, nil
}

func (u Users) DeleteUserByGoogleId(connection *gorm.DB, googleUserId string) error {
	user := &Users{}
	result := connection.First(&user, "google_user_id = ?", googleUserId)
	if result.RowsAffected == 0 {
		return ErrRecordNotFound
	}

	connection.Delete(&user)

	return nil
}
