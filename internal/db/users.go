package db

type Users struct {
	ID           int    `json:"id" gorm:"primaryKey"`
	GoogleUserId string `json:"google_user_id"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}
