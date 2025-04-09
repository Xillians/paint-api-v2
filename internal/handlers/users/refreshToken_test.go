package users_test

import (
	"context"
	"paint-api/internal/handlers/users"
	"paint-api/internal/jwt"
	"paint-api/internal/middleware"
	"paint-api/internal/testutils"
	"testing"
)

func TestRefreshTokenHandler(t *testing.T) {
	connection, cleanUp := testutils.OpenTestConnection()
	userId := testData.User.GoogleUserId
	t.Run("Refresh token with valid userId", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, middleware.DbKey, connection)
		ctx = context.WithValue(ctx, middleware.UserIdKey, userId)
		jwtService := jwt.NewJWTService("some_secret")
		ctx = context.WithValue(ctx, middleware.JwtKey, *jwtService)

		registerInput := &users.RefreshTokenInput{}
		output, err := users.RefreshTokenHandler(ctx, registerInput)
		if err != nil {
			t.Fatalf("Failed to register user: %v", err)
		}
		if output.Body.Token == "" {
			t.Errorf("Expected non-empty token, got empty")
		}
		if output.Body.ExpiresAt == "" {
			t.Errorf("Expected non-empty expires_at, got empty")
		}
	})
	t.Run("Refresh token with missing userId", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, middleware.DbKey, connection)
		jwtService := jwt.NewJWTService("some_secret")
		ctx = context.WithValue(ctx, middleware.JwtKey, *jwtService)

		registerInput := &users.RefreshTokenInput{}
		output, err := users.RefreshTokenHandler(ctx, registerInput)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
		if output != nil {
			t.Errorf("Expected nil output, got %v", output)
		}

	})
	t.Run("Refresh token without jwt service", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, middleware.DbKey, connection)
		ctx = context.WithValue(ctx, middleware.UserIdKey, userId)

		registerInput := &users.RefreshTokenInput{}
		output, err := users.RefreshTokenHandler(ctx, registerInput)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
		if output != nil {
			t.Errorf("Expected nil output, got %v", output)
		}
	})
	t.Run("Refresh token without db connection", func(t *testing.T) {
		ctx := context.Background()
		jwtService := jwt.NewJWTService("some_secret")
		ctx = context.WithValue(ctx, middleware.JwtKey, *jwtService)
		ctx = context.WithValue(ctx, middleware.UserIdKey, userId)

		registerInput := &users.RefreshTokenInput{}
		output, err := users.RefreshTokenHandler(ctx, registerInput)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
		if output != nil {
			t.Errorf("Expected nil output, got %v", output)
		}
	})
	t.Run("Refresh token with closed db connection", func(t *testing.T) {
		ctx, err := createClosedDBContext()
		if err != nil {
			t.Fatalf("Failed to create closed DB context: %v", err)
		}
		ctx = context.WithValue(ctx, middleware.UserIdKey, userId)
		jwtService := jwt.NewJWTService("some_secret")
		ctx = context.WithValue(ctx, middleware.JwtKey, *jwtService)

		registerInput := &users.RefreshTokenInput{}
		output, err := users.RefreshTokenHandler(ctx, registerInput)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
		if output != nil {
			t.Errorf("Expected nil output, got %v", output)
		}
	})
	t.Cleanup(cleanUp)
}
