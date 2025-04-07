package brands_test

import (
	"context"
	"paint-api/internal/handlers/brands"
	"paint-api/internal/middleware"
	"paint-api/internal/testutils"
	"testing"
)

func TestUpdateHandler(t *testing.T) {
	connection, cleanUp := testutils.OpenTestConnection()
	t.Run("Update brand with valid data", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, middleware.DbKey, connection)
		ctx = context.WithValue(ctx, middleware.RoleKey, "administrator")

		input := &brands.UpdateBrandInput{
			ID: uint(testData.Brand.ID),
			Body: brands.UpdateBrandInputBody{
				Name: "Updated Brand",
			},
		}
		output, err := brands.UpdateHandler(ctx, input)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}

		if output.Body.ID != int(input.ID) {
			t.Errorf("Expected ID %d, got %d", input.ID, output.Body.ID)
		}
		if output.Body.Name != input.Body.Name {
			t.Errorf("Expected name %s, got %s", input.Body.Name, output.Body.Name)
		}
		if output.Body.UpdatedAt == testData.Brand.UpdatedAt {
			t.Errorf("Expected updated_at to be different, got %s", output.Body.UpdatedAt)
		}
	})
	t.Run("Update with invalid role", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, middleware.DbKey, connection)
		ctx = context.WithValue(ctx, middleware.RoleKey, "user")
		input := &brands.UpdateBrandInput{
			ID: uint(testData.Brand.ID),
			Body: brands.UpdateBrandInputBody{
				Name: "Updated Brand",
			},
		}
		_, err := brands.UpdateHandler(ctx, input)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
		if err.Error() != "You are not allowed to perform this action" {
			t.Errorf("Expected forbidden error, got %v", err.Error())
		}
	})
	t.Run("Update with invalid ID", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, middleware.DbKey, connection)
		ctx = context.WithValue(ctx, middleware.RoleKey, "administrator")

		input := &brands.UpdateBrandInput{
			ID: uint(999999),
			Body: brands.UpdateBrandInputBody{
				Name: "Updated Brand",
			},
		}
		_, err := brands.UpdateHandler(ctx, input)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
		if err.Error() != "Brand not found" {
			t.Errorf("Expected not found error, got %v", err.Error())
		}
	})
	t.Run("Update without db connection", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, middleware.RoleKey, "administrator")

		input := &brands.UpdateBrandInput{
			ID: uint(testData.Brand.ID),
			Body: brands.UpdateBrandInputBody{
				Name: "Updated Brand",
			},
		}

		_, err := brands.UpdateHandler(ctx, input)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
		if err.Error() != "failed to update brand" {
			t.Errorf("Expected database connection error, got %v", err.Error())
		}
	})
	t.Run("Update with closed db connection", func(t *testing.T) {
		conn, _ := testutils.OpenTestConnection()

		ctx := context.Background()
		ctx = context.WithValue(ctx, middleware.DbKey, conn)
		ctx = context.WithValue(ctx, middleware.RoleKey, "administrator")

		db, err := conn.DB()
		if err != nil {
			t.Fatalf("Failed to get DB from connection: %v", err)
		}
		db.Close()

		input := &brands.UpdateBrandInput{
			ID: uint(testData.Brand.ID),
			Body: brands.UpdateBrandInputBody{
				Name: "Updated Brand",
			},
		}

		_, err = brands.UpdateHandler(ctx, input)
		if err == nil {
			t.Errorf("Expected error, got nil")
		}
		if err.Error() != "Failed to update brand" {
			t.Errorf("Expected database connection error, got %v", err.Error())
		}
	})
	t.Cleanup(cleanUp)
}
