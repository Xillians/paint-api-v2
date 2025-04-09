package users_test

import (
	"context"
	"paint-api/internal/db"
	"paint-api/internal/handlers/users"
	"paint-api/internal/middleware"
	"paint-api/internal/testutils"
	"testing"
)

func TestRegisterUserHandler(t *testing.T) {
	connection, cleanUp := testutils.OpenTestConnection()
	t.Run("Delete a user with valid data", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, middleware.DbKey, connection)

		registerInput := &users.RegisterUserInput{
			Body: db.RegisterUserInput{
				GoogleUserId: "123454321",
				Email:        "asd@dsa.io",
			},
		}
		_, err := users.RegisterHandler(ctx, registerInput)
		if err != nil {
			t.Fatalf("Failed to register user: %v", err)
		}

		ctx = context.WithValue(ctx, middleware.UserIdKey, "123454321")

		input := &users.ForgetUserInput{}
		output, err := users.ForgetHandler(ctx, input)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if output.Body != "User deleted successfully" {
			t.Errorf("Expected 'User deleted successfully', got %v", output.Body)
		}
	})
	t.Run("Delete with missing userId", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, middleware.DbKey, connection)

		input := &users.ForgetUserInput{}
		output, err := users.ForgetHandler(ctx, input)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
		if output != nil {
			t.Errorf("Expected nil output, got %v", output)
		}
	})
	t.Run("Delete with missing db context", func(t *testing.T) {
		ctx := context.Background()
		input := &users.ForgetUserInput{}
		output, err := users.ForgetHandler(ctx, input)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
		if output != nil {
			t.Errorf("Expected nil output, got %v", output)
		}
	})
	t.Run("db connection error", func(t *testing.T) {
		ctx, err := createClosedDBContext()
		if err != nil {
			t.Fatalf("Failed to create closed DB context: %v", err)
		}

		ctx = context.WithValue(ctx, middleware.UserIdKey, "123454321")

		input := &users.ForgetUserInput{}

		output, err := users.ForgetHandler(ctx, input)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
		if output != nil {
			t.Errorf("Expected nil output, got %v", output)
		}
	})
	t.Cleanup(cleanUp)
}
