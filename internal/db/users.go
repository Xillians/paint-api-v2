package db

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
