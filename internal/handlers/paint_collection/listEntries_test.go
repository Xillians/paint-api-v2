package paint_collection_test

import (
	"context"
	"paint-api/internal/db"
	"paint-api/internal/handlers/paint_collection"
	"paint-api/internal/handlers/users"
	"paint-api/internal/middleware"
	"paint-api/internal/testutils"
	"testing"
)

func TestListEntriesHandler(t *testing.T) {
	connection, cleanUp := testutils.OpenTestConnection()
	testUser := testData.User
	t.Run("Add entry and find entry", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, middleware.DbKey, connection)
		ctx = context.WithValue(ctx, middleware.UserIdKey, testUser.GoogleUserId)

		createEntryInput := &paint_collection.AddToCollectionInput{
			Body: paint_collection.AddToCollectionInputBody{
				Quantity: 1,
				PaintID:  testData.Paint.Id,
			},
		}
		createEntryOutput, err := paint_collection.CreateHandler(ctx, createEntryInput)
		if err != nil {
			t.Fatalf("Failed to create entry: %v", err)
		}

		listEntriesInput := &paint_collection.ListPaintCollectionInput{}
		listEntriesOutput, err := paint_collection.ListHandler(ctx, listEntriesInput)
		if err != nil {
			t.Fatalf("Failed to list entries: %v", err)
		}
		mappedEntries := mapEntries(listEntriesOutput.Body.Collection)
		if len(mappedEntries) == 0 {
			t.Fatalf("Expected non-empty collection, got empty")
		}
		entry, exists := mappedEntries[createEntryOutput.Body.ID]
		if !exists {
			t.Fatalf("Expected entry with ID %d, got none", testData.Paint.Id)
		}
		if entry.Quantity != createEntryInput.Body.Quantity {
			t.Fatalf("Expected quantity %d, got %d", createEntryInput.Body.Quantity, entry.Quantity)
		}
	})
	t.Run("Get empty list when user not in db", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, middleware.DbKey, connection)
		ctx = context.WithValue(ctx, middleware.UserIdKey, "2")

		input := &paint_collection.ListPaintCollectionInput{}
		output, err := paint_collection.ListHandler(ctx, input)
		if err != nil {
			t.Fatalf("Failed to list entries: %v", err)
		}

		mappedEntries := mapEntries(output.Body.Collection)
		if len(mappedEntries) != 0 {
			t.Fatalf("Expected non-empty collection, got empty")
		}
	})
	t.Run("New user only seees their own paints (none)", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, middleware.DbKey, connection)

		userInput := &users.RegisterUserInput{
			Body: db.RegisterUserInput{
				GoogleUserId: "4201",
				Email:        testUser.Email,
			},
		}
		newUser, err := users.RegisterHandler(ctx, userInput)
		if err != nil {
			t.Fatalf("Failed to register user: %v", err)
		}
		ctx = context.WithValue(ctx, middleware.UserIdKey, newUser.Body.GoogleUserId)

		input := &paint_collection.ListPaintCollectionInput{}
		output, err := paint_collection.ListHandler(ctx, input)
		if err != nil {
			t.Fatalf("Failed to list entries: %v", err)
		}
		mappedEntries := mapEntries(output.Body.Collection)
		if len(mappedEntries) != 0 {
			t.Fatalf("Expected empty collection, got %d entries", len(mappedEntries))
		}

		deleteUserInput := &users.ForgetUserInput{}
		_, err = users.ForgetHandler(ctx, deleteUserInput)
		if err != nil {
			t.Fatalf("Failed to delete user: %v", err)
		}

	})
	t.Run("Missing DB context", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, middleware.UserIdKey, testUser.GoogleUserId)

		input := &paint_collection.ListPaintCollectionInput{}
		_, err := paint_collection.ListHandler(ctx, input)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})
	t.Run("Missing userId context", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, middleware.DbKey, connection)

		input := &paint_collection.ListPaintCollectionInput{}
		_, err := paint_collection.ListHandler(ctx, input)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})
	t.Run("Closed DB connection", func(t *testing.T) {
		ctx, err := createClosedDBContext()
		if err != nil {
			t.Fatalf("Failed to create closed DB context: %v", err)
		}
		ctx = context.WithValue(ctx, middleware.UserIdKey, testUser.GoogleUserId)

		input := &paint_collection.ListPaintCollectionInput{}
		_, err = paint_collection.ListHandler(ctx, input)
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})
	t.Cleanup(cleanUp)
}

func mapEntries(entries []db.CollectionPaintDetails) map[int]db.CollectionPaintDetails {
	mappedEntries := make(map[int]db.CollectionPaintDetails)
	for _, entry := range entries {
		mappedEntries[entry.ID] = entry
	}
	return mappedEntries
}
