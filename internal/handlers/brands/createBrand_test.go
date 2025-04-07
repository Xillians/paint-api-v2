package brands_test

import (
	"context"
	"paint-api/internal/db"
	"paint-api/internal/handlers/brands"
	"paint-api/internal/middleware"
	"paint-api/internal/testutils"
	"testing"
)

func TestCreateHandler(t *testing.T) {
	connection, cleanUp := testutils.OpenTestConnection()
	t.Run("Create brand with valid data", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, middleware.DbKey, connection)
		input := &brands.CreatebrandInput{
			Body: db.CreateBrandInput{
				Name: "Test Brand",
			},
		}
		output, err := brands.CreateHandler(ctx, input)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if output == nil {
			t.Errorf("Expected output, got nil")
		}
	})
	t.Run("Create with missing db context", func(t *testing.T) {
		ctx := context.Background()
		input := &brands.CreatebrandInput{
			Body: db.CreateBrandInput{
				Name: "Test Brand",
			},
		}
		output, err := brands.CreateHandler(ctx, input)
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

		input := &brands.CreatebrandInput{
			Body: db.CreateBrandInput{
				Name: "Test Brand",
			},
		}

		output, err := brands.CreateHandler(ctx, input)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
		if output != nil {
			t.Errorf("Expected nil output, got %v", output)
		}
	})
	t.Cleanup(cleanUp)
}
