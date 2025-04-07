package brands_test

import (
	"context"
	"paint-api/internal/handlers/brands"
	"paint-api/internal/middleware"
	"paint-api/internal/testutils"
	"testing"
)

func TestGetbrandHelper(t *testing.T) {
	connection, cleanUp := testutils.OpenTestConnection()
	t.Run("Get brand with valid data", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, middleware.DbKey, connection)
		input := &brands.GetBrandInput{
			ID: uint(testData.Brand.ID),
		}

		output, err := brands.GetHandler(ctx, input)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if uint(output.Body.ID) != uint(testData.Brand.ID) {
			t.Errorf("Expected brand ID %d, got %d", testData.Brand.ID, output.Body.ID)
		}
	})
	t.Run("Get brand with invalid ID", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, middleware.DbKey, connection)
		input := &brands.GetBrandInput{
			ID: 99999,
		}

		output, err := brands.GetHandler(ctx, input)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
		if output != nil {
			t.Errorf("Expected nil output, got %v", output)
		}
	})
	t.Run("Get brand with missing db context", func(t *testing.T) {
		ctx := context.Background()
		input := &brands.GetBrandInput{
			ID: uint(testData.Brand.ID),
		}

		output, err := brands.GetHandler(ctx, input)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
		if output != nil {
			t.Errorf("Expected nil output, got %v", output)
		}
	})
	t.Run("Get brand with closed db connection", func(t *testing.T) {
		ctx, err := createClosedDBContext()
		if err != nil {
			t.Fatalf("Failed to create closed DB context: %v", err)
		}
		input := &brands.GetBrandInput{
			ID: uint(testData.Brand.ID),
		}

		output, err := brands.GetHandler(ctx, input)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
		if output != nil {
			t.Errorf("Expected nil output, got %v", output)
		}
	})
	t.Cleanup(cleanUp)
}
