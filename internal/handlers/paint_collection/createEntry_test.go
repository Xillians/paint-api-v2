package paint_collection_test

import (
	"context"
	"paint-api/internal/handlers/paint_collection"
	"paint-api/internal/middleware"
	"paint-api/internal/testutils"
	"testing"
)

func TestCreateEntryHandler(t *testing.T) {
	connection, cleanUp := testutils.OpenTestConnection()
	testUser := testData.User
	defaultInput := &paint_collection.AddToCollectionInput{
		Body: paint_collection.AddToCollectionInputBody{
			Quantity: 1,
			PaintID:  testData.Paint.Id,
		},
	}
	t.Run("Successfully add an entry", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, middleware.DbKey, connection)
		ctx = context.WithValue(ctx, middleware.UserIdKey, testUser.GoogleUserId)

		createEntryOutput, err := paint_collection.CreateHandler(ctx, defaultInput)
		if err != nil {
			t.Fatalf("Failed to create entry: %v", err)
		}
		if createEntryOutput.Body.ID == 0 {
			t.Fatalf("Expected non-zero ID, got %d", createEntryOutput.Body.ID)
		}
		if createEntryOutput.Body.Quantity != defaultInput.Body.Quantity {
			t.Fatalf("Expected quantity %d, got %d", defaultInput.Body.Quantity, createEntryOutput.Body.Quantity)
		}
		if createEntryOutput.Body.PaintID != defaultInput.Body.PaintID {
			t.Fatalf("Expected PaintID %d, got %d", defaultInput.Body.PaintID, createEntryOutput.Body.PaintID)
		}
		if createEntryOutput.Body.User.GoogleUserId != testUser.GoogleUserId {
			t.Fatalf("Expected UserID %s, got %s", testUser.GoogleUserId, createEntryOutput.Body.User.GoogleUserId)
		}
		if createEntryOutput.Body.Paint.Brand.ID != testData.Paint.Brand.ID {
			t.Fatalf("Expected BrandID %d, got %d", testData.Paint.Brand.ID, createEntryOutput.Body.Paint.Brand.ID)
		}

		deleteInput := &paint_collection.DeleteCollectionEntryInput{
			Id: createEntryOutput.Body.ID,
		}
		_, err = paint_collection.DeleteHandler(ctx, deleteInput)
		if err != nil {
			t.Fatalf("Failed to delete entry: %v", err)
		}

	})
	t.Run("Fail to find userId in db", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, middleware.DbKey, connection)
		ctx = context.WithValue(ctx, middleware.UserIdKey, "1001")

		_, err := paint_collection.CreateHandler(ctx, defaultInput)
		if err == nil {
			t.Fatalf("Expected error, got none")
		}
		if err.Error() != "could not add paint to collection" {
			t.Fatalf("Expected error 'could not add paint to collection', got '%v'", err)
		}
	})
	t.Run("Fail to find db in context", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, middleware.UserIdKey, testUser.GoogleUserId)

		_, err := paint_collection.CreateHandler(ctx, defaultInput)
		if err == nil {
			t.Fatalf("Expected error, got none")
		}
		if err.Error() != "could not add paint to collection" {
			t.Fatalf("Expected error 'could not add paint to collection', got '%v'", err)
		}
	})
	t.Run("Fail to find userId in context", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, middleware.DbKey, connection)

		_, err := paint_collection.CreateHandler(ctx, defaultInput)
		if err == nil {
			t.Fatalf("Expected error, got none")
		}
		if err.Error() != "could not add paint to collection" {
			t.Fatalf("Expected error 'could not add paint to collection', got '%v'", err)
		}
	})
	t.Run("DB error", func(t *testing.T) {
		ctx, err := createClosedDBContext()
		if err != nil {
			t.Fatalf("Failed to create closed DB context: %v", err)
		}
		ctx = context.WithValue(ctx, middleware.UserIdKey, testUser.GoogleUserId)

		_, err = paint_collection.CreateHandler(ctx, defaultInput)
		if err == nil {
			t.Fatalf("Expected error, got none")
		}
		if err.Error() != "could not add paint to collection" {
			t.Fatalf("Expected error 'could not add paint to collection', got '%v'", err)
		}
	})
	t.Cleanup(cleanUp)
}
