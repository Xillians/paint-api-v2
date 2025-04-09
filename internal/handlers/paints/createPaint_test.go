package paints_test

import (
	"context"
	"paint-api/internal/db"
	"paint-api/internal/handlers/paints"
	"paint-api/internal/middleware"
	"paint-api/internal/testutils"
	"testing"
)

func TestCreatePaintHanlder(t *testing.T) {
	connection, cleanUp := testutils.OpenTestConnection()
	brand := testData.Brand
	t.Run("Successfully create paint", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, middleware.DbKey, connection)

		input := &paints.CreatePaintInput{
			Body: db.CreatePaintInput{
				Name:      "Test Paint",
				ColorCode: "#FFFFFF",
				BrandId:   brand.ID,
			},
		}
		output, err := paints.CreateHandler(ctx, input)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if output.Body.Name != input.Body.Name {
			t.Fatalf("Expected paint name %s, got %s", input.Body.Name, output.Body.Name)
		}
		if output.Body.ColorCode != input.Body.ColorCode {
			t.Fatalf("Expected color code %s, got %s", input.Body.ColorCode, output.Body.ColorCode)
		}
		if output.Body.BrandId != brand.ID {
			t.Fatalf("Expected brand ID %d, got %d", brand.ID, output.Body.BrandId)
		}
	})
	t.Run("Validate correct hex color", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, middleware.DbKey, connection)

		input := &paints.CreatePaintInput{
			Body: db.CreatePaintInput{
				Name:      "Test Paint",
				ColorCode: "#FFFFFF",
				BrandId:   brand.ID,
			},
		}
		output, err := paints.CreateHandler(ctx, input)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if output.Body.ColorCode != input.Body.ColorCode {
			t.Fatalf("Expected color code %s, got %s", input.Body.ColorCode, output.Body.ColorCode)
		}
		if output.Body.BrandId != brand.ID {
			t.Fatalf("Expected brand ID %d, got %d", brand.ID, output.Body.BrandId)
		}
	})
	t.Run("Validate incorrect hex color", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, middleware.DbKey, connection)

		input := &paints.CreatePaintInput{
			Body: db.CreatePaintInput{
				Name:      "Test Paint",
				ColorCode: "FFFFFF",
				BrandId:   brand.ID,
			},
		}
		output, err := paints.CreateHandler(ctx, input)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
		if output != nil {
			t.Fatalf("Expected nil output, got %v", output)
		}
	})
	t.Run("Missing DB context", func(t *testing.T) {
		ctx := context.Background()

		input := &paints.CreatePaintInput{
			Body: db.CreatePaintInput{
				Name:      "Test Paint",
				ColorCode: "#FFFFFF",
				BrandId:   brand.ID,
			},
		}
		output, err := paints.CreateHandler(ctx, input)
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

		input := &paints.CreatePaintInput{
			Body: db.CreatePaintInput{
				Name:      "Test Paint",
				ColorCode: "#FFFFFF",
				BrandId:   brand.ID,
			},
		}
		output, err := paints.CreateHandler(ctx, input)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
		if output != nil {
			t.Fatalf("Expected nil output, got %v", output)
		}
	})
	t.Cleanup(cleanUp)
}
