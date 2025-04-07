package users_test

import (
	"net/http"
	"testing"
)

func TestRegisterUser(t *testing.T) {
	t.Run("Successfully create a user", func(t *testing.T) {
		response, err := createTestUser("12345", "asd@asd.io")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if response.Result().StatusCode != http.StatusOK {
			t.Fatalf("Expected status code 200, got %d", response.Result().StatusCode)
		}
	})
	t.Run("Attempt to create duplicate user", func(t *testing.T) {
		response, err := createTestUser("1234", "asd@asd.io")
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if response.Result().StatusCode != http.StatusOK {
			t.Fatalf("Expected status code 200, got %d", response.Result().StatusCode)
		}

		registerUserInput := map[string]interface{}{
			"email":   "asd@asd.io",
			"user_id": "1234",
		}
		createResponse := testApi.Post("/register", registerUserInput)

		if createResponse.Result().StatusCode != http.StatusConflict {
			t.Fatalf("Expected status code 409, got %d", createResponse.Result().StatusCode)
		}
	})
	t.Run("Attempt to create user with invalid email", func(t *testing.T) {
		registerUserInput := map[string]interface{}{
			"email":   "invalid-email",
			"user_id": "1234",
		}
		createResponse := testApi.Post("/register", registerUserInput)
		if createResponse.Result().StatusCode != http.StatusBadRequest {
			t.Fatalf("Expected status code 400, got %d", createResponse.Result().StatusCode)
		}
	})
	t.Run("Attempt to create user with invalid user_id", func(t *testing.T) {

		registerUserInput := map[string]interface{}{
			"email":   "asd@asd.io",
			"user_id": 1,
		}
		createResponse := testApi.Post("/register", registerUserInput)
		if createResponse.Result().StatusCode != http.StatusUnprocessableEntity {
			t.Fatalf("Expected status code 422, got %d", createResponse.Result().StatusCode)
		}
	})
}
