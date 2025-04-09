package paints_test

import (
	"context"
	"paint-api/internal/db"
	"paint-api/internal/handlers/paints"
	"paint-api/internal/middleware"
	"paint-api/internal/testutils"
	"testing"
)

func TestUpdateHandler(t *testing.T) {
	connection, cleanUp := testutils.OpenTestConnection()
	brand := testData.Brand
	colorCode := "#FFFFFF"
	paintName := "Test Paint"
	description := "Test description"
	updatedPaintName := "Updated Paint"

	t.Run("Successfully update paint", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, middleware.DbKey, connection)
		ctx = context.WithValue(ctx, middleware.RoleKey, "administrator")

		createInput := &paints.CreatePaintInput{
			Body: db.CreatePaintInput{
				Name:        paintName,
				ColorCode:   colorCode,
				BrandId:     brand.ID,
				Description: description,
			},
		}
		createOutput, err := paints.CreateHandler(ctx, createInput)
		if err != nil {
			t.Fatalf("Failed to create paint: %v", err)
		}
		paintId := createOutput.Body.Id

		updateInput := &paints.UpdatePaintInput{
			Id: paintId,
			Body: paints.UpdatePaintInputBody{
				Name: updatedPaintName,
			},
		}
		updateOutput, err := paints.UpdateHandler(ctx, updateInput)
		if err != nil {
			t.Fatalf("Failed to update paint: %v", err)
		}
		if updateOutput.Body.Name != updatedPaintName {
			t.Errorf("Expected updated paint name '%s', got '%s'", updatedPaintName, updateOutput.Body.Name)
		}
	})
	t.Run("Fail to update paint with non-existent ID", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, middleware.DbKey, connection)
		ctx = context.WithValue(ctx, middleware.RoleKey, "administrator")

		updateInput := &paints.UpdatePaintInput{
			Id: 99999, //
			Body: paints.UpdatePaintInputBody{
				Name: updatedPaintName,
			},
		}

		_, err := paints.UpdateHandler(ctx, updateInput)
		if err == nil {
			t.Fatalf("Expected error when updating non-existent paint, got nil")
		}
		if err.Error() != "paint not found" {
			t.Errorf("Expected 'paint not found' error, got '%v'", err)
		}
	})
	t.Run("Fail to update paint without administrator role", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, middleware.DbKey, connection)
		ctx = context.WithValue(ctx, middleware.RoleKey, "user")

		input := &paints.UpdatePaintInput{
			Id: 1,
			Body: paints.UpdatePaintInputBody{
				Name: updatedPaintName,
			},
		}

		_, err := paints.UpdateHandler(ctx, input)
		if err == nil {
			t.Fatalf("Expected error when updating paint without administrator role, got nil")
		}
		if err.Error() != "You are not allowed to perform this action" {
			t.Errorf("Expected 'You are not allowed to perform this action' error, got '%v'", err)
		}
	})
	t.Run("fail to update with missing authentication", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, middleware.DbKey, connection)

		input := &paints.UpdatePaintInput{
			Id: 1,
			Body: paints.UpdatePaintInputBody{
				Name: updatedPaintName,
			},
		}

		_, err := paints.UpdateHandler(ctx, input)
		if err == nil {
			t.Fatalf("Expected error when updating paint without authentication, got nil")
		}
		if err.Error() != "failed to update paint" {
			t.Errorf("Expected 'failed to update paint' error, got '%v'", err)
		}
	})
	t.Run("Fail with missing db context", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, middleware.RoleKey, "administrator")
		input := &paints.UpdatePaintInput{
			Id: 1,
			Body: paints.UpdatePaintInputBody{
				Name: updatedPaintName,
			},
		}

		_, err := paints.UpdateHandler(ctx, input)
		if err == nil {
			t.Fatalf("Expected error when updating paint with missing db context, got nil")
		}
		if err.Error() != "failed to update paint" {
			t.Fatalf("Expected 'failed to update paint' error, got '%v'", err)
		}
	})
	t.Run("Fail to update paint with closed db connection", func(t *testing.T) {
		ctx, err := createClosedDBContext()
		if err != nil {
			t.Fatalf("Failed to create closed DB context: %v", err)
		}
		ctx = context.WithValue(ctx, middleware.RoleKey, "administrator")

		input := &paints.UpdatePaintInput{
			Id: 1,
			Body: paints.UpdatePaintInputBody{
				Name: updatedPaintName,
			},
		}

		_, err = paints.UpdateHandler(ctx, input)
		if err == nil {
			t.Fatalf("Expected error when updating paint with closed DB connection, got nil")
		}
		if err.Error() != "failed to update paint" {
			t.Fatalf("Expected 'failed to update paint' error, got '%v'", err)
		}
	})
	t.Cleanup(cleanUp)
}
