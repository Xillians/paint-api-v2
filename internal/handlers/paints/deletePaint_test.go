package paints_test

import (
	"context"
	"paint-api/internal/db"
	"paint-api/internal/handlers/paints"
	"paint-api/internal/middleware"
	"paint-api/internal/testutils"
	"testing"
)

func TestDeletePaintHandler(t *testing.T) {
	connection, cleanUp := testutils.OpenTestConnection()
	testBrand := testData.Brand
	colorCode := "#FFFFFF"
	paintName := "Test Paint"
	description := "Test description"

	t.Run("Successfully delete paint", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, middleware.DbKey, connection)
		ctx = context.WithValue(ctx, middleware.RoleKey, "administrator")

		createInput := &paints.CreatePaintInput{
			Body: db.CreatePaintInput{
				Name:        paintName,
				ColorCode:   colorCode,
				BrandId:     testBrand.ID,
				Description: description,
			},
		}
		createOutput, err := paints.CreateHandler(ctx, createInput)
		if err != nil {
			t.Fatalf("Failed to create paint: %v", err)
		}
		paintId := createOutput.Body.Id

		deleteInput := &paints.DeletePaintInput{
			Id: paintId,
		}
		deleteOutput, err := paints.DeleteHandler(ctx, deleteInput)
		if err != nil {
			t.Fatalf("Failed to delete paint: %v", err)
		}
		if deleteOutput.Body != "Paint deleted successfully" {
			t.Errorf("Expected 'Paint deleted successfully', got '%s'", deleteOutput.Body)
		}
	})
	t.Run("Fail to delete paint with non-existent ID", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, middleware.DbKey, connection)
		ctx = context.WithValue(ctx, middleware.RoleKey, "administrator")

		deleteInput := &paints.DeletePaintInput{
			Id: 99999,
		}

		_, err := paints.DeleteHandler(ctx, deleteInput)
		if err == nil {
			t.Fatalf("Expected error when deleting non-existent paint, got nil")
		}
	})
	t.Run("Fail to delete paint with non-administrator role", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, middleware.DbKey, connection)
		ctx = context.WithValue(ctx, middleware.RoleKey, "user") // Non-administrator role

		createInput := &paints.CreatePaintInput{
			Body: db.CreatePaintInput{
				Name:        paintName,
				ColorCode:   colorCode,
				BrandId:     testBrand.ID,
				Description: description,
			},
		}

		createOutput, err := paints.CreateHandler(ctx, createInput)
		if err != nil {
			t.Fatalf("Failed to create paint: %v", err)
		}
		paintId := createOutput.Body.Id
		deleteInput := &paints.DeletePaintInput{
			Id: paintId,
		}

		_, err = paints.DeleteHandler(ctx, deleteInput)
		if err == nil {
			t.Fatalf("Expected error when deleting paint with non-administrator role, got nil")
		}
		expectedError := "You are not allowed to perform this action"
		if err.Error() != expectedError {
			t.Fatalf("Expected error message '%s', got '%s'", expectedError, err.Error())
		}
	})
	t.Run("Fail to delete with missing authentication", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, middleware.DbKey, connection)

		deleteInput := &paints.DeletePaintInput{
			Id: 1,
		}
		_, err := paints.DeleteHandler(ctx, deleteInput)
		if err == nil {
			t.Fatalf("Expected error when authentication is missing, got nil")
		}
		expectedError := "Failed to delete paint"
		if err.Error() != expectedError {
			t.Fatalf("Expected error message '%s', got '%s'", expectedError, err.Error())
		}
	})
	t.Run("Missing db connection", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, middleware.RoleKey, "administrator")

		deleteInput := &paints.DeletePaintInput{
			Id: 1,
		}
		_, err := paints.DeleteHandler(ctx, deleteInput)
		if err == nil {
			t.Fatalf("Expected error when db connection is missing, got nil")
		}
	})
	t.Cleanup(cleanUp)
}
