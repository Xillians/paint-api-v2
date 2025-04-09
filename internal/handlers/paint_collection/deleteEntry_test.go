package paint_collection_test

import (
	"context"
	"paint-api/internal/handlers/paint_collection"
	"paint-api/internal/middleware"
	"paint-api/internal/testutils"
	"testing"
)

func TestDeleteEntryHandler(t *testing.T) {
	connection, cleanUp := testutils.OpenTestConnection()
	testUser := testData.User
	defaultInput := &paint_collection.AddToCollectionInput{
		Body: paint_collection.AddToCollectionInputBody{
			Quantity: 1,
			PaintID:  testData.Paint.Id,
		},
	}
	t.Run("Successfully delete entry", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, middleware.DbKey, connection)
		ctx = context.WithValue(ctx, middleware.UserIdKey, testUser.GoogleUserId)

		createEntryOutput, err := paint_collection.CreateHandler(ctx, defaultInput)
		if err != nil {
			t.Fatalf("Failed to create entry: %v", err)
		}

		deleteInput := &paint_collection.DeleteCollectionEntryInput{
			Id: createEntryOutput.Body.ID,
		}
		output, err := paint_collection.DeleteHandler(ctx, deleteInput)
		if err != nil {
			t.Fatalf("Failed to delete entry: %v", err)
		}
		if output.Body != "Entry deleted successfully" {
			t.Fatalf("Expected success message, got %v", output.Body)
		}
	})
	t.Run("Fail to delete non-existent entry", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, middleware.DbKey, connection)
		ctx = context.WithValue(ctx, middleware.UserIdKey, testUser.GoogleUserId)

		deleteInput := &paint_collection.DeleteCollectionEntryInput{
			Id: 9999,
		}

		_, err := paint_collection.DeleteHandler(ctx, deleteInput)
		if err == nil {
			t.Fatalf("Expected error, got none")
		}
	})
	t.Run("Fail to find db in context", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, middleware.UserIdKey, testUser.GoogleUserId)

		deleteInput := &paint_collection.DeleteCollectionEntryInput{
			Id: 1,
		}
		_, err := paint_collection.DeleteHandler(ctx, deleteInput)
		if err == nil {
			t.Fatalf("Expected error, got none")
		}
		if err.Error() != "failed to delete entry" {
			t.Fatalf("Expected error 'failed to delete entry', got '%v'", err)
		}
	})
	t.Run("Fail to find userId in context", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, middleware.DbKey, connection)

		deleteInput := &paint_collection.DeleteCollectionEntryInput{
			Id: 1,
		}
		_, err := paint_collection.DeleteHandler(ctx, deleteInput)
		if err == nil {
			t.Fatalf("Expected error, got none")
		}
		if err.Error() != "Entry not found" {
			t.Fatalf("Expected error 'Entry not found', got '%v'", err)
		}
	})
	t.Run("DB error", func(t *testing.T) {
		ctx, err := createClosedDBContext()
		if err != nil {
			t.Fatalf("Failed to create closed DB context: %v", err)
		}
		ctx = context.WithValue(ctx, middleware.UserIdKey, testUser.GoogleUserId)

		deleteInput := &paint_collection.DeleteCollectionEntryInput{
			Id: 1,
		}
		_, err = paint_collection.DeleteHandler(ctx, deleteInput)
		if err == nil {
			t.Fatalf("Expected error, got none")
		}
		if err.Error() != "Entry not found" {
			t.Fatalf("Expected error 'failed to delete entry', got '%v'", err)
		}
	})
	t.Cleanup(cleanUp)
}
