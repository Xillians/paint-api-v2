package users_test

import (
	"context"
	"paint-api/internal/db"
	"paint-api/internal/handlers/users"
	"paint-api/internal/middleware"
	"paint-api/internal/testutils"
	"testing"
)

func TestRegisterHandler(t *testing.T) {
	connection, cleanUp := testutils.OpenTestConnection()
	userId := testData.User.GoogleUserId
	createdUserId := "4201"

	t.Run("Register new user", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, middleware.DbKey, connection)

		registerInput := &users.RegisterUserInput{
			Body: db.RegisterUserInput{
				GoogleUserId: createdUserId,
				Email:        "asd@dsa.io",
			},
		}
		output, err := users.RegisterHandler(ctx, registerInput)
		if err != nil {
			t.Fatalf("Failed to register user: %v", err)
		}
		if output.Body.CreatedAt == "" {
			t.Errorf("Expected non-empty userId, got empty")
		}
		if output.Body.UpdatedAt == "" {
			t.Errorf("Expected non-empty userId, got empty")
		}
		if output.Body.GoogleUserId != createdUserId {
			t.Errorf("Expected userId %s, got %s", createdUserId, output.Body.GoogleUserId)
		}
		if output.Body.Email == "" {
			t.Errorf("Expected non-empty email, got empty")
		}
		if output.Body.Role != "user" {
			t.Errorf("Expected role 'user', got %s", output.Body.Role)
		}

		// delete user
		deleteInput := &users.ForgetUserInput{}
		ctx = context.WithValue(ctx, middleware.UserIdKey, createdUserId)
		deleteOutput, err := users.ForgetHandler(ctx, deleteInput)
		if err != nil {
			t.Fatalf("Failed to delete user: %v", err)
		}
		if deleteOutput == nil {
			t.Errorf("Expected non-nil delete output, got nil")
		}
	})
	t.Run("Register existing user", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, middleware.DbKey, connection)

		registerInput := &users.RegisterUserInput{
			Body: db.RegisterUserInput{
				GoogleUserId: userId,
				Email:        "asd@dsa.io",
			},
		}
		output, err := users.RegisterHandler(ctx, registerInput)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
		if output != nil {
			t.Errorf("Expected nil output, got %v", output)
		}
	})
	t.Run("Register with invalid email", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, middleware.DbKey, connection)

		registerInput := &users.RegisterUserInput{
			Body: db.RegisterUserInput{
				GoogleUserId: createdUserId,
				Email:        "invalid-email",
			},
		}
		output, err := users.RegisterHandler(ctx, registerInput)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
		if output != nil {
			t.Errorf("Expected nil output, got %v", output)
		}
	})
	t.Run("Register without db context", func(t *testing.T) {
		ctx := context.Background()
		registerInput := &users.RegisterUserInput{
			Body: db.RegisterUserInput{
				GoogleUserId: createdUserId,
				Email:        "invalid-email",
			},
		}
		output, err := users.RegisterHandler(ctx, registerInput)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
		if output != nil {
			t.Errorf("Expected nil output, got %v", output)
		}
	})
	t.Run("Register with closed db connection", func(t *testing.T) {
		ctx, err := createClosedDBContext()
		if err != nil {
			t.Fatalf("Failed to create closed DB context: %v", err)
		}
		registerInput := &users.RegisterUserInput{
			Body: db.RegisterUserInput{
				GoogleUserId: createdUserId,
				Email:        "invalid-email",
			},
		}
		output, err := users.RegisterHandler(ctx, registerInput)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
		if output != nil {
			t.Errorf("Expected nil output, got %v", output)
		}
	})
	t.Cleanup(cleanUp)
}
