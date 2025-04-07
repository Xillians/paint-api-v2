package brands_test

import (
	"context"
	"paint-api/internal/handlers/brands"
	"paint-api/internal/middleware"
	"paint-api/internal/testutils"
	"testing"
)

func TestListBrandHelper(t *testing.T) {
	connection, cleanUp := testutils.OpenTestConnection()
	t.Run("List brand with valid data", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, middleware.DbKey, connection)

		input := &brands.ListBrandInput{}
		output, err := brands.ListHandler(ctx, input)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if output == nil {
			t.Errorf("Expected output, got nil")
		}
	})
	t.Run("List brand with missing db context", func(t *testing.T) {
		ctx := context.Background()
		input := &brands.ListBrandInput{}

		output, err := brands.ListHandler(ctx, input)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
		if output != nil {
			t.Errorf("Expected nil output, got %v", output)
		}
	})
	t.Run("List brand with closed db connection", func(t *testing.T) {
		ctx, err := createClosedDBContext()
		if err != nil {
			t.Fatalf("Failed to create closed DB context: %v", err)
		}
		input := &brands.ListBrandInput{}

		output, err := brands.ListHandler(ctx, input)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
		if output != nil {
			t.Errorf("Expected nil output, got %v", output)
		}
	})
	t.Cleanup(cleanUp)
}
