package brands_test

import (
	"context"
	"paint-api/internal/db"
	"paint-api/internal/handlers/brands"
	"paint-api/internal/middleware"
	"paint-api/internal/testutils"
	"testing"
)

func TestDeleteHandler(t *testing.T) {
	connection, cleanUp := testutils.OpenTestConnection()
	t.Run("Delete brand with valid data", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, middleware.DbKey, connection)
		ctx = context.WithValue(ctx, middleware.RoleKey, "administrator")

		// Create a brand to delete
		brandInput := &brands.CreatebrandInput{
			Body: db.CreateBrandInput{
				Name: "Test Brand",
			},
		}
		createOutput, err := brands.CreateHandler(ctx, brandInput)
		if err != nil {
			t.Fatalf("Failed to create brand: %v", err)
		}

		input := &brands.DeleteBrandInput{
			ID: uint(createOutput.Body.ID),
		}
		output, err := brands.DeleteHandler(ctx, input)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if output == nil {
			t.Errorf("Expected output, got nil")
		}
	})
	t.Run("Delete with invalid role", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, middleware.DbKey, connection)
		ctx = context.WithValue(ctx, middleware.RoleKey, "user")

		input := &brands.DeleteBrandInput{
			ID: 1,
		}

		_, err := brands.DeleteHandler(ctx, input)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
		error := err.Error()
		if error != "You are not allowed to perform this action" {
			t.Errorf("Expected forbidden error, got %v", error)
		}
	})
	t.Run("Attempt to delete non-existent brand", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, middleware.DbKey, connection)
		ctx = context.WithValue(ctx, middleware.RoleKey, "administrator")

		input := &brands.DeleteBrandInput{
			ID: 9999,
		}

		output, err := brands.DeleteHandler(ctx, input)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}

		if output != nil {
			t.Errorf("Expected nil output, got %v", output)
		}
	})
	t.Run("Delete with missing db context", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, middleware.RoleKey, "administrator")

		input := &brands.DeleteBrandInput{
			ID: 1,
		}

		output, err := brands.DeleteHandler(ctx, input)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}

		if output != nil {
			t.Errorf("Expected nil output, got %v", output)
		}
	})
	t.Run("db connection error", func(t *testing.T) {
		connection, _ := testutils.OpenTestConnection()
		ctx := context.Background()
		ctx = context.WithValue(ctx, middleware.RoleKey, "administrator")

		sql, err := connection.DB()
		if err != nil {
			t.Fatalf("Failed to get DB from connection: %v", err)
		}

		sql.Close()

		ctx = context.WithValue(ctx, middleware.DbKey, connection)

		input := &brands.DeleteBrandInput{
			ID: 1,
		}

		output, err := brands.DeleteHandler(ctx, input)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}

		if output != nil {
			t.Errorf("Expected nil output, got %v", output)
		}
	})
	t.Cleanup(cleanUp)
}
