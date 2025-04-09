package paints_test

import (
	"context"
	"paint-api/internal/db"
	"paint-api/internal/handlers/paints"
	"paint-api/internal/middleware"
	"paint-api/internal/testutils"
	"testing"
)

func TestListPaintHandler(t *testing.T) {
	connection, cleanUp := testutils.OpenTestConnection()
	brand := testData.Brand
	paint := testData.Paint
	t.Run("Successfully list paints", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, middleware.DbKey, connection)

		input := &paints.ListPaintInput{}
		output, err := paints.ListHandler(ctx, input)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		paintMap := make(map[int]db.Paints)
		for _, p := range output.Body.Paints {
			paintMap[p.Id] = p
		}
		if len(paintMap) == 0 {
			t.Fatalf("Expected to find at least one paint, got none")
		}
		if paintMap[paint.Id].Id != paint.Id {
			t.Fatalf("Expected paint ID %d, got %d", paint.Id, paintMap[paint.Id].Id)
		}
		if paintMap[paint.Id].BrandId != brand.ID {
			t.Fatalf("Expected to find brand data %d, got %d", brand.ID, paintMap[paint.Id].BrandId)
		}
	})
	t.Run("Missing DB context", func(t *testing.T) {
		ctx := context.Background()

		input := &paints.ListPaintInput{}
		output, err := paints.ListHandler(ctx, input)
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

		input := &paints.ListPaintInput{}
		output, err := paints.ListHandler(ctx, input)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
		if output != nil {
			t.Fatalf("Expected nil output, got %v", output)
		}
	})
	t.Cleanup(cleanUp)
}
