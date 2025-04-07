package users_test

import (
	"net/http"
	"paint-api/internal/handlers/users"
	"testing"
)

func TestForgetUser(t *testing.T) {
	deleteUserEndpoint := "/forget"
	t.Run("Create and delete user", func(t *testing.T) {
		createResponse := createTestUser("123454321", "asd@ghj.io")
		if createResponse.Result().StatusCode != http.StatusOK {
			t.Fatalf("Expected status code 200, got %d", createResponse.Result().StatusCode)
		}

		loginResponse := loginTestUser("123454321")
		if loginResponse.Result().StatusCode != http.StatusOK {
			t.Fatalf("Expected status code 200, got %d", loginResponse.Result().StatusCode)
		}

		loginResponseBody, err := parseResponse[users.LoginOutputBody](loginResponse)
		if err != nil {
			t.Fatalf("Failed to decode login response: %v", err)
		}

		response := deleteTestUser(loginResponseBody.Token)
		if response.Result().StatusCode != http.StatusOK {
			t.Fatalf("Expected status code 200, got %d", response.Result().StatusCode)
		}
	})
	t.Run("Delete user with invalid token", func(t *testing.T) {
		deleteUserHeader := makeRequestHeader("invalid_token")
		deleteResponse := testApi.Delete(deleteUserEndpoint, deleteUserHeader)
		if deleteResponse.Result().StatusCode != http.StatusUnauthorized {
			t.Fatalf("Expected status code 401, got %d", deleteResponse.Result().StatusCode)
		}
	})
	t.Run("Attempt to double delete user", func(t *testing.T) {
		createUserResponse := createTestUser("123454321", "asd@ghj.io")
		if createUserResponse.Result().StatusCode != http.StatusOK {
			t.Fatalf("Expected status code 200, got %d", createUserResponse.Result().StatusCode)
		}

		loginResponse := loginTestUser("123454321")
		if loginResponse.Result().StatusCode != http.StatusOK {
			t.Fatalf("Expected status code 200, got %d", loginResponse.Result().StatusCode)
		}

		loginResponseBody, err := parseResponse[users.LoginOutputBody](loginResponse)
		if err != nil {
			t.Fatalf("Failed to decode login response: %v", err)
		}

		deleteUserResponse := deleteTestUser(loginResponseBody.Token)
		if deleteUserResponse.Result().StatusCode != http.StatusOK {
			t.Fatalf("Expected status code 200, got %d", deleteUserResponse.Result().StatusCode)
		}

		deleteUserResponse = deleteTestUser(loginResponseBody.Token)
		if deleteUserResponse.Result().StatusCode != http.StatusNotFound {
			t.Fatalf("Expected status code 401, got %d", deleteUserResponse.Result().StatusCode)
		}
	})
}
