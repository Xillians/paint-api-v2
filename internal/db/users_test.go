package db_test

import (
	"paint-api/internal/db"
	"testing"
)

func createTestUser() *db.Users {
	input := db.RegisterUserInput{
		GoogleUserId: "1234567890",
		Email:        "test@testerson.io",
	}
	user, err := db.Users{}.RegisterUser(testDB, input)
	if err != nil {
		return nil
	}
	return user
}

func TestCreateUser(t *testing.T) {
	testInput := db.RegisterUserInput{
		GoogleUserId: "123423",
		Email:        "",
	}
	t.Run("Register user", func(t *testing.T) {
		user, err := db.Users{}.RegisterUser(testDB, testInput)
		if err != nil {
			t.Errorf("Error registering user: %v", err)
		}

		err = db.Users{}.DeleteUserByGoogleId(testDB, user.GoogleUserId)
		if err != nil {
			t.Errorf("Error deleting user by google id: %v", err)
		}
	})
	t.Run("Attempt to register user with existing google user id", func(t *testing.T) {
		_, err := db.Users{}.RegisterUser(testDB, testInput)
		if err != nil {
			t.Errorf("Error registering user: %v", err)
		}

		_, err = db.Users{}.RegisterUser(testDB, testInput)
		if err == nil {
			t.Errorf("Expected duplicate error!")
		}
	})
	t.Run("Transaction failure", func(t *testing.T) {
		connection := OpenTestConnection()
		sql, _ := connection.DB()
		sql.Close()

		_, err := db.Users{}.RegisterUser(connection, testInput)
		if err == nil {
			t.Errorf("Expected error registering user with nil db, got nil")
		}
	})
}

func TestGetUserByGoogleId(t *testing.T) {
	user := createTestUser()
	t.Run("Get user by google id", func(t *testing.T) {
		result, err := db.Users{}.GetUserByGoogleId(testDB, user.GoogleUserId)
		if err != nil {
			t.Errorf("Error getting user by google id: %v", err)
		}
		if result == nil {
			t.Errorf("Expected user, got nil")
		}
	})
	t.Run("Attempt to get user by non-existent google id", func(t *testing.T) {
		_, err := db.Users{}.GetUserByGoogleId(testDB, "001")
		if err == nil {
			t.Errorf("Expected error getting user by non-existent google id, got nil")
		}
	})
	t.Run("Transaction failure", func(t *testing.T) {
		connection := OpenTestConnection()
		sql, _ := connection.DB()
		sql.Close()

		_, err := db.Users{}.GetUserByGoogleId(connection, user.GoogleUserId)
		if err == nil {
			t.Errorf("Expected error getting user by google id with nil db, got nil")
		}
	})
	t.Cleanup(func() {
		err := db.Users{}.DeleteUserByGoogleId(testDB, user.GoogleUserId)
		if err != nil {
			t.Errorf("Error deleting user by google id: %v", err)
		}
	})
}

func TestDeleteUser(t *testing.T) {
	user := createTestUser()
	t.Run("Delete user by google id", func(t *testing.T) {

		err := db.Users{}.DeleteUserByGoogleId(testDB, user.GoogleUserId)
		if err != nil {
			t.Errorf("Error deleting user by google id: %v", err)
		}
	})
	t.Run("Attempt to delete user by non-existent google id", func(t *testing.T) {
		err := db.Users{}.DeleteUserByGoogleId(testDB, "001")
		if err == nil {
			t.Errorf("Expected error deleting user by non-existent google id, got nil")
		}
	})
	t.Run("Transaction failure", func(t *testing.T) {
		user := createTestUser()
		connection := OpenTestConnection()
		sql, _ := connection.DB()
		sql.Close()

		err := db.Users{}.DeleteUserByGoogleId(connection, user.GoogleUserId)
		if err == nil {
			t.Errorf("Expected error deleting user by google id with nil db, got nil")
		}
	})
	t.Cleanup(func() {
		err := db.Users{}.DeleteUserByGoogleId(testDB, user.GoogleUserId)
		if err != nil {
			t.Errorf("Error deleting user by google id: %v", err)
		}
	})
}
