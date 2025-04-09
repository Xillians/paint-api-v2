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

func TestUpdateEntryHandler(t *testing.T) {
	connection, cleanUp := testutils.OpenTestConnection()
	testUser := testData.User
	testPaint := testData.Paint
	defaultCreateInput := paint_collection.AddToCollectionInput{
		Body: paint_collection.AddToCollectionInputBody{
			PaintID:  testPaint.Id,
			Quantity: 1,
		},
	}
	t.Run("Successfully update an entry", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, middleware.DbKey, connection)
		ctx = context.WithValue(ctx, middleware.UserIdKey, testUser.GoogleUserId)

		// Create a new collection entry
		createOutput, err := paint_collection.CreateHandler(ctx, &defaultCreateInput)
		if err != nil {
			t.Fatalf("Failed to create collection entry: %v", err)
		}

		updateInput := paint_collection.UpdateCollectionEntryInput{
			Id: createOutput.Body.ID,
			Body: paint_collection.UpdateCollectionEntryInputBody{
				Quantity: 2,
				PaintId:  testPaint.Id,
			},
		}
		updateOutput, err := paint_collection.UpdateHandler(ctx, &updateInput)
		if err != nil {
			t.Fatalf("Failed to update collection entry: %v", err)
		}
		if updateOutput.Body.Quantity != 2 {
			t.Errorf("Expected quantity to be 2, got %d", updateOutput.Body.Quantity)
		}
		if updateOutput.Body.PaintID != testPaint.Id {
			t.Errorf("Expected paint ID to be %d, got %d", testPaint.Id, updateOutput.Body.PaintID)
		}
	})
	t.Run("only owner can update entry", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, middleware.DbKey, connection)
		ctx = context.WithValue(ctx, middleware.UserIdKey, testUser.GoogleUserId)

		// Create a new collection entry
		createOutput, err := paint_collection.CreateHandler(ctx, &defaultCreateInput)
		if err != nil {
			t.Fatalf("Failed to create collection entry: %v", err)
		}

		// create new user
		newUserInput := users.RegisterUserInput{
			Body: db.RegisterUserInput{
				GoogleUserId: "4201",
				Email:        "asd@dsa.io",
			},
		}
		newUserOutput, err := users.RegisterHandler(ctx, &newUserInput)
		if err != nil {
			t.Fatalf("Failed to create new user: %v", err)
		}
		ctx = context.WithValue(ctx, middleware.UserIdKey, newUserOutput.Body.GoogleUserId)

		updateInput := paint_collection.UpdateCollectionEntryInput{
			Id: createOutput.Body.ID,
			Body: paint_collection.UpdateCollectionEntryInputBody{
				Quantity: 2,
				PaintId:  testPaint.Id,
			},
		}
		_, err = paint_collection.UpdateHandler(ctx, &updateInput)
		if err == nil {
			t.Fatalf("Expected error when updating entry as non-owner, got nil")
		}
		if err.Error() != "entry not found" {
			t.Fatalf("Expected 'entry not found' error, got: %v", err)
		}

		deleteInput := users.ForgetUserInput{}
		_, err = users.ForgetHandler(ctx, &deleteInput)
		if err != nil {
			t.Fatalf("Failed to delete new user: %v", err)
		}
	})
	t.Run("missing db context", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, middleware.UserIdKey, testUser.GoogleUserId)

		updateInput := paint_collection.UpdateCollectionEntryInput{
			Id: 1,
			Body: paint_collection.UpdateCollectionEntryInputBody{
				Quantity: 2,
				PaintId:  testPaint.Id,
			},
		}
		_, err := paint_collection.UpdateHandler(ctx, &updateInput)
		if err == nil {
			t.Fatalf("Expected error when missing db context, got nil")
		}
		if err.Error() != "Failed to update entry." {
			t.Fatalf("Expected 'Failed to update entry.' error, got: %v", err)
		}
	})
	t.Run("missing user context", func(t *testing.T) {
		ctx := context.Background()
		ctx = context.WithValue(ctx, middleware.DbKey, connection)

		updateInput := paint_collection.UpdateCollectionEntryInput{
			Id: 1,
			Body: paint_collection.UpdateCollectionEntryInputBody{
				Quantity: 2,
				PaintId:  testPaint.Id,
			},
		}
		_, err := paint_collection.UpdateHandler(ctx, &updateInput)
		if err == nil {
			t.Fatalf("Expected error when missing user context, got nil")
		}
		if err.Error() != "entry not found" {
			t.Fatalf("Expected 'entry not found' error, got: %v", err)
		}
	})
	t.Run("DB connection closed", func(t *testing.T) {
		ctx, err := createClosedDBContext()
		if err != nil {
			t.Fatalf("Failed to create closed DB context: %v", err)
		}
		ctx = context.WithValue(ctx, middleware.UserIdKey, testUser.GoogleUserId)
		updateInput := paint_collection.UpdateCollectionEntryInput{
			Id: 1,
			Body: paint_collection.UpdateCollectionEntryInputBody{
				Quantity: 2,
				PaintId:  testPaint.Id,
			},
		}

		_, err = paint_collection.UpdateHandler(ctx, &updateInput)
		if err == nil {
			t.Fatalf("Expected error when DB connection is closed, got nil")
		}
		if err.Error() != "entry not found" {
			t.Fatalf("Expected 'entry not found' error, got: %v", err)
		}
	})
	t.Cleanup(cleanUp)
}
