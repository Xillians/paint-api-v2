package db_test

import (
	"errors"
	"paint-api/internal/db"
	"testing"
)

func createTestEntry() (*db.CollectionPaintDetails, error) {
	paint := createTestPaint()
	user := createTestUser()
	entryInput := db.CreateCollectionEntryInput{
		PaintID:  paint.Id,
		Quantity: 1,
		UserId:   user.ID,
	}
	entry, err := db.CollectionPaintDetails{}.CreateEntry(testDB, entryInput)
	if err != nil {
		return nil, err
	}
	return entry, nil
}

func TestCreateEntry(t *testing.T) {
	user := createTestUser()
	paint := createTestPaint()
	t.Run("Create entry", func(t *testing.T) {
		entryInput := db.CreateCollectionEntryInput{
			PaintID:  paint.Id,
			Quantity: 1,
			UserId:   user.ID,
		}
		entry, err := db.CollectionPaintDetails{}.CreateEntry(testDB, entryInput)
		if err != nil {
			t.Errorf("Error creating entry: %v", err)
		}

		err = db.CollectionPaintDetails{}.DeleteEntry(testDB, entry.ID)
		if err != nil {
			t.Errorf("Error deleting entry by id: %v", err)
		}
	})
	t.Run("Transaction error", func(t *testing.T) {
		connection := OpenTestConnection()
		sql, err := connection.DB()
		if err != nil {
			t.Errorf("Error getting sql connection: %v", err)
		}
		sql.Close()

		entryInput := db.CreateCollectionEntryInput{
			PaintID:  paint.Id,
			Quantity: 1,
			UserId:   user.ID,
		}
		_, err = db.CollectionPaintDetails{}.CreateEntry(connection, entryInput)
		if err == nil {
			t.Error("Expected error creating entry")
		}
	})
	t.Cleanup(func() {
		err := db.Users{}.DeleteUserByGoogleId(testDB, user.GoogleUserId)
		if err != nil {
			t.Errorf("Error deleting user by google id: %v", err)
		}
		err = db.Paints{}.DeletePaint(testDB, paint.Id)
		if err != nil {
			t.Errorf("Error deleting paint by id: %v", err)
		}
	})
}

func TestGetEntry(t *testing.T) {
	entry, err := createTestEntry()
	if err != nil {
		t.Errorf("Error creating test entry: %v", err)
	}

	t.Run("Get entry", func(t *testing.T) {
		_, err := db.CollectionPaintDetails{}.GetEntry(testDB, entry.ID, entry.User.GoogleUserId)
		if err != nil {
			t.Errorf("Error getting entry: %v", err)
		}
	})
	t.Run("Attempt to get non-existent entry", func(t *testing.T) {
		_, err := db.CollectionPaintDetails{}.GetEntry(testDB, 0, entry.User.GoogleUserId)
		if err == nil {
			t.Error("Expected error getting non-existent entry")
		}
		if !errors.Is(err, db.ErrRecordNotFound) {
			t.Errorf("Wrong error getting non-existent entry: %v", err)
		}
	})
	t.Run("Attempt to get existing entry with non-existent user", func(t *testing.T) {
		_, err := db.CollectionPaintDetails{}.GetEntry(testDB, entry.ID, "001")
		if err == nil {
			t.Error("Expected error getting entry with non-existent user")
		}
		if !errors.Is(err, db.ErrRecordNotFound) {
			t.Errorf("Wrong error getting entry with non-existent user: %v", err)
		}
	})
	t.Run("Attempt to access entry with wrong user", func(t *testing.T) {
		newUserInput := db.RegisterUserInput{
			GoogleUserId: "new_user",
			Email:        "asd@fgh.io",
		}
		newUser, err := db.Users{}.RegisterUser(testDB, newUserInput, "user")
		if err != nil {
			t.Errorf("Error registering new user: %v", err)
		}
		_, err = db.CollectionPaintDetails{}.GetEntry(testDB, entry.ID, newUser.GoogleUserId)
		if err == nil {
			t.Error("Expected error getting entry with wrong user")
		}
		if !errors.Is(err, db.ErrRecordNotFound) {
			t.Errorf("Wrong error getting entry with wrong user: %v", err)
		}
	})

	t.Cleanup(func() {
		cleanUp(entry, t)
	})
}

func TestListEntries(t *testing.T) {
	entry, err := createTestEntry()
	if err != nil {
		t.Errorf("Error creating test entry: %v", err)
	}

	t.Run("List entries", func(t *testing.T) {
		_, err := db.CollectionPaintDetails{}.ListEntries(testDB, entry.User.GoogleUserId)
		if err != nil {
			t.Errorf("Error listing entries: %v", err)
		}
	})
	t.Run("Attempt to list entries with non-existent user", func(t *testing.T) {
		entries, err := db.CollectionPaintDetails{}.ListEntries(testDB, "001")
		if err != nil {
			t.Errorf("Error listing entries with non-existent user: %v", err)
		}
		if len(entries) > 0 {
			t.Errorf("Expected no entries, got %v", len(entries))
		}
	})
	t.Run("Transaction error", func(t *testing.T) {
		connection := OpenTestConnection()
		sql, err := connection.DB()
		if err != nil {
			t.Errorf("Error getting sql connection: %v", err)
		}
		sql.Close()

		_, err = db.CollectionPaintDetails{}.ListEntries(connection, entry.User.GoogleUserId)
		if err == nil {
			t.Error("Expected error listing entries")
		}
	})
	t.Cleanup(func() {
		cleanUp(entry, t)
	})
}

func TestUpdateEntry(t *testing.T) {
	entry, err := createTestEntry()
	if err != nil {
		t.Errorf("Error creating test entry: %v", err)
	}

	t.Run("Update entry", func(t *testing.T) {
		updateInput := db.UpdateCollectionEntryInput{
			ID:       entry.ID,
			PaintID:  entry.Paint.Id,
			Quantity: 2,
		}
		_, err := db.CollectionPaintDetails{}.UpdateEntry(testDB, updateInput)
		if err != nil {
			t.Errorf("Error updating entry: %v", err)
		}
	})
	t.Run("Attempt to update non-existent entry", func(t *testing.T) {
		updateInput := db.UpdateCollectionEntryInput{
			ID:       0,
			PaintID:  entry.Paint.Id,
			Quantity: 2,
		}
		_, err := db.CollectionPaintDetails{}.UpdateEntry(testDB, updateInput)
		if err == nil {
			t.Error("Expected error updating non-existent entry")
		}
		if !errors.Is(err, db.ErrRecordNotFound) {
			t.Errorf("Wrong error updating non-existent entry: %v", err)
		}
	})
	t.Run("Transaction error", func(t *testing.T) {
		connection := OpenTestConnection()
		sql, err := connection.DB()
		if err != nil {
			t.Errorf("Error getting sql connection: %v", err)
		}
		sql.Close()

		updateInput := db.UpdateCollectionEntryInput{
			ID:       entry.ID,
			PaintID:  entry.Paint.Id,
			Quantity: 2,
		}
		_, err = db.CollectionPaintDetails{}.UpdateEntry(connection, updateInput)
		if err == nil {
			t.Error("Expected error updating entry")
		}
	})

	t.Cleanup(func() {
		cleanUp(entry, t)
	})
}

func cleanUp(entry *db.CollectionPaintDetails, t *testing.T) {

	err := db.CollectionPaintDetails{}.DeleteEntry(testDB, entry.ID)
	if err != nil {
		t.Errorf("Error deleting entry by id: %v", err)
	}
	err = db.Users{}.DeleteUserByGoogleId(testDB, entry.User.GoogleUserId)
	if err != nil {
		t.Errorf("Error deleting user by google id: %v", err)
	}
	err = db.Paints{}.DeletePaint(testDB, entry.Paint.Id)
	if err != nil {
		t.Errorf("Error deleting paint by id: %v", err)
	}
}
