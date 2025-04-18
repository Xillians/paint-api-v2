package paints_test

import (
	"context"
	"paint-api/internal/handlers/paints"
	"paint-api/internal/middleware"
	"paint-api/internal/testutils"
	"testing"
)

func TestGetPaintHandler(t *testing.T) {
	connection, cleanUp := testutils.OpenTestConnection()
	brand := testData.Brand
	paint := testData.Paint
	t.Run("Successfully get paint", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, middleware.DbKey, connection)

		input := &paints.GetPaintsInput{
			Id: paint.Id,
		}
		output, err := paints.GetHandler(ctx, input)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if output.Body.Id != paint.Id {
			t.Fatalf("Expected paint ID %d, got %d", paint.Id, output.Body.Id)
		}
		if output.Body.BrandId != brand.ID {
			t.Fatalf("Expected brand ID %d, got %d", brand.ID, output.Body.BrandId)
		}
	})
	t.Run("Try to get non-existing paint", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, middleware.DbKey, connection)

		input := &paints.GetPaintsInput{
			Id: 99999,
		}
		output, err := paints.GetHandler(ctx, input)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
		if output != nil {
			t.Fatalf("Expected nil output, got %v", output)
		}
	})
	t.Run("Missing DB context", func(t *testing.T) {
		ctx := context.Background()

		input := &paints.GetPaintsInput{
			Id: paint.Id,
		}
		output, err := paints.GetHandler(ctx, input)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
		if output != nil {
			t.Fatalf("Expected nil output, got %v", output)
		}
	})
	t.Run("DB connection error", func(t *testing.T) {
		ctx, err := createClosedDBContext()
		if err != nil {
			t.Fatalf("Failed to create closed DB context: %v", err)
		}

		input := &paints.GetPaintsInput{
			Id: paint.Id,
		}
		output, err := paints.GetHandler(ctx, input)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
		if output != nil {
			t.Fatalf("Expected nil output, got %v", output)
		}
	})
	t.Cleanup(cleanUp)
}
