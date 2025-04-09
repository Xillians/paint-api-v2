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

	t.Run("Register new user", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, middleware.DbKey, connection)

		registerInput := &users.RegisterUserInput{
			Body: db.RegisterUserInput{
				GoogleUserId: "123454321",
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
		if output.Body.GoogleUserId == "" {
			t.Errorf("Expected non-empty userId, got empty")
		}
		if output.Body.Email == "" {
			t.Errorf("Expected non-empty email, got empty")
		}
		if output.Body.Role != "user" {
			t.Errorf("Expected role 'user', got %s", output.Body.Role)
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
				GoogleUserId: "4202",
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
				GoogleUserId: "4321",
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
				GoogleUserId: "1211",
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
