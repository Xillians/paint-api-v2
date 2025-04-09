package users_test

import (
	"context"
	"paint-api/internal/handlers/users"
	"paint-api/internal/jwt"
	"paint-api/internal/middleware"
	"paint-api/internal/testutils"
	"testing"
)

func TestLoginHandler(t *testing.T) {
	connection, cleanUp := testutils.OpenTestConnection()
	testUserId := testData.User.GoogleUserId

	t.Run("Login with valid userId", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, middleware.DbKey, connection)
		jwtService := jwt.NewJWTService("some_secret")
		ctx = context.WithValue(ctx, middleware.JwtKey, *jwtService)

		input := &users.LoginInput{
			GoogleUserId: testUserId,
		}
		output, err := users.LoginHandler(ctx, input)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if output.Body.Token == "" {
			t.Errorf("Expected non-empty token, got empty")
		}
		if output.Body.ExpiresAt == "" {
			t.Errorf("Expected non-empty expires_at, got empty")
		}
	})
	t.Run("Login with missing userId", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, middleware.DbKey, connection)
		jwtService := jwt.NewJWTService("some_secret")
		ctx = context.WithValue(ctx, middleware.JwtKey, *jwtService)

		input := &users.LoginInput{
			GoogleUserId: "123",
		}
		output, err := users.LoginHandler(ctx, input)
		if err == nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if output != nil {
			t.Errorf("Expected nil output, got %v", output)
		}
	})
	t.Run("Login without db connection", func(t *testing.T) {
		ctx := context.Background()
		jwtService := jwt.NewJWTService("some_secret")
		ctx = context.WithValue(ctx, middleware.JwtKey, *jwtService)

		input := &users.LoginInput{
			GoogleUserId: testUserId,
		}
		output, err := users.LoginHandler(ctx, input)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
		if output != nil {
			t.Errorf("Expected nil output, got %v", output)
		}
	})
	t.Run("Login without jwt service", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, middleware.DbKey, connection)

		input := &users.LoginInput{
			GoogleUserId: testUserId,
		}
		output, err := users.LoginHandler(ctx, input)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
		if output != nil {
			t.Errorf("Expected nil output, got %v", output)
		}
	})
	t.Run("Login with closed db connection", func(t *testing.T) {
		ctx, err := createClosedDBContext()
		if err != nil {
			t.Errorf("Failed to create closed DB context: %v", err)
		}

		jwtService := jwt.NewJWTService("some_secret")
		ctx = context.WithValue(ctx, middleware.JwtKey, *jwtService)

		input := &users.LoginInput{
			GoogleUserId: testUserId,
		}

		output, err := users.LoginHandler(ctx, input)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
		if output != nil {
			t.Errorf("Expected nil output, got %v", output)
		}
	})
	t.Cleanup(cleanUp)
}
