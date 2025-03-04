package handlers

import (
	"context"
	"errors"
	"net/http"
	"paint-api/internal/db"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

type createUserInputBody struct {
	GoogleUserId string `json:"google_user_id" validate:"required"`
}
type createUserInput struct {
	Body createUserInputBody
}

type createUserOutput struct {
	Body db.Users
}

var CreateUserOperation = huma.Operation{
	Method: http.MethodPost,
	Path:   "/Users",
	Tags:   []string{"Users"},
}

func CreateUserHandler(ctx context.Context, input *createUserInput) (*createUserOutput, error) {
	connection, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return nil, errors.New("could not retrieve db from context")
	}
	User := db.Users{
		GoogleUserId: input.Body.GoogleUserId,
		CreatedAt:    time.Now().String(),
		UpdatedAt:    time.Now().String(),
	}
	connection.Create(&User)
	return &createUserOutput{Body: User}, nil
}

type listUserInput struct {
}

type listUserOutput struct {
	Body []db.Users
}

var ListUsersOperation = huma.Operation{
	Method: http.MethodGet,
	Path:   "/Users",
	Tags:   []string{"Users"},
}

func ListUsersHandler(ctx context.Context, input *listUserInput) (*listUserOutput, error) {
	connection, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return nil, errors.New("could not retrieve db from context")
	}
	var Users []db.Users
	connection.Find(&Users)
	return &listUserOutput{Body: Users}, nil
}

type getUsersInput struct {
	Id int `path:"id"`
}

type getUserOutput struct {
	Body db.Users
}

var GetUsersOperation = huma.Operation{
	Method: http.MethodGet,
	Path:   "/Users/{id}",
	Tags:   []string{"Users"},
}

func GetUserHandler(ctx context.Context, input *getUsersInput) (*getUserOutput, error) {
	connection, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return nil, errors.New("could not retrieve db from context")
	}
	var User db.Users
	if err := connection.First(&User, input.Id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, huma.NewError(http.StatusNotFound, "User not found")
		}
		return nil, err
	}
	return &getUserOutput{Body: User}, nil
}

type deleteUserInput struct {
	Id int `path:"id"`
}

type deleteUserOutput struct {
	Body string
}

var DeleteUserOperation = huma.Operation{
	Method: http.MethodDelete,
	Path:   "/Users/{id}",
	Tags:   []string{"Users"},
}

func DeleteUserHandler(ctx context.Context, input *deleteUserInput) (*deleteUserOutput, error) {
	connection, ok := ctx.Value("db").(*gorm.DB)
	if !ok {
		return nil, errors.New("could not retrieve db from context")
	}
	if err := connection.Delete(&db.Users{}, input.Id).Error; err != nil {
		return nil, err
	}
	return &deleteUserOutput{Body: "User deleted successfully"}, nil
}
