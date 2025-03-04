package handlers

import (
	"context"
	"errors"
	"net/http"
	"paint-api/internal/db"
	"paint-api/internal/jwt"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

type registerUserInputBody struct {
	GoogleUserId string `json:"user_id" validate:"required"`
}
type RegisterUserInput struct {
	Body registerUserInputBody
}

type registerUserOutput struct {
	Body db.Users
}

var RegisterUserOperation = huma.Operation{
	Method: http.MethodPost,
	Path:   "/register",
	Tags:   []string{"Users"},
}

func RegisterUserHandler(ctx context.Context, input *RegisterUserInput) (*registerUserOutput, error) {
	connection, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return nil, errors.New("could not retrieve db from context")
	}

	user := connection.First(&db.Users{}, "google_user_id = ?", input.Body.GoogleUserId)
	if user.RowsAffected > 0 {
		return nil, huma.NewError(http.StatusConflict, "User already exists")
	}

	User := db.Users{
		GoogleUserId: input.Body.GoogleUserId,
		CreatedAt:    time.Now().String(),
		UpdatedAt:    time.Now().String(),
	}
	connection.Create(&User)

	return &registerUserOutput{Body: User}, nil
}

type LoginInput struct {
	GoogleUserId int `path:"id"`
}

type LoginOutputBody struct {
	Token string `json:"token"`
}

type LoginOutput struct {
	Body LoginOutputBody
}

var LoginOperation = huma.Operation{
	Method: http.MethodGet,
	Path:   "/login/{id}",
	Tags:   []string{"Users"},
}

func LoginHandler(ctx context.Context, input *LoginInput) (*LoginOutput, error) {
	connection, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return nil, errors.New("could not retrieve db from context")
	}
	jwt, ok := ctx.Value("jwtKey").(jwt.JWTService)
	if !ok {
		return nil, errors.New("could not retrieve jwt from context")
	}

	var User db.Users
	if err := connection.First(&User, "google_user_id = ?", input.GoogleUserId).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, huma.NewError(http.StatusNotFound, "User not found")
		}
		return nil, err
	}
	// Since we have found the user, we want to create a new session for the user
	// This session should be a JWT that will be used to authenticate the user
	// All other requests will require this JWT to be passed in the header
	// The JWT will be generated using the user's google user id
	// The JWT will be signed using a secret key that will be stored
	// in the environment variables
	// The JWT will have an expiry time of 1 hour
	token, err := jwt.GenerateToken(User.GoogleUserId)
	if err != nil {
		return nil, err
	}
	return &LoginOutput{Body: LoginOutputBody{Token: token}}, nil
}

type refreshTokenInput struct {
	Id int `path:"id"`
}

type refreshTokenOutputBody struct {
	Token string `json:"token"`
}

type refreshTokenOutput struct {
	Body refreshTokenOutputBody
}

var RefreshTokenOperation = huma.Operation{
	Method: http.MethodGet,
	Path:   "/refresh/{id}",
	Tags:   []string{"Users"},
}

func RefreshTokenHandler(ctx context.Context, input *refreshTokenInput) (*refreshTokenOutput, error) {
	// We will use the google user id to generate a new JWT
	// This JWT will have an expiry time of 1 hour
	// The JWT will be signed using a secret key that will be stored
	// in the environment variables
	return nil, nil
}

type forgetUserInput struct {
	Id int `path:"id"`
}

type forgetUserOutput struct {
	Body string
}

var ForgetUserOperation = huma.Operation{
	Method: http.MethodDelete,
	Path:   "/forget/{id}",
	Tags:   []string{"Users"},
}

func ForgetUserHandler(ctx context.Context, input *forgetUserInput) (*forgetUserOutput, error) {
	connection, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return nil, errors.New("could not retrieve db from context")
	}

	user := connection.First(&db.Users{}, "google_user_id = ?", input.Id)
	if user.RowsAffected == 0 {
		return nil, huma.NewError(http.StatusNotFound, "User not found")
	}

	if err := connection.Delete(&db.Users{}, "google_user_id = ?", input.Id).Error; err != nil {
		return nil, err
	}
	return &forgetUserOutput{Body: "User deleted successfully"}, nil
}
