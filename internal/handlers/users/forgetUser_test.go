package users_test

import (
	"net/http"
	"paint-api/internal/handlers/users"
	"testing"
)

func TestForgetUser(t *testing.T) {
	deleteUserEndpoint := "/forget"
	t.Run("Create and delete user", func(t *testing.T) {
		_, err := createTestUser("123454321", "asd@ghj.io")
		if err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}

		loginResponse, err := loginTestUser("123454321")
		if err != nil {
			t.Fatalf("Failed to login test user: %v", err)
		}

		loginResponseBody, err := parseResponse[users.LoginOutputBody](loginResponse)
		if err != nil {
			t.Fatalf("Failed to decode login response: %v", err)
		}

		deleteUserHeader := makeRequestHeader(loginResponseBody.Token)
		deleteResponse := testApi.Delete(deleteUserEndpoint, deleteUserHeader)
		if deleteResponse.Result().StatusCode != http.StatusOK {
			t.Fatalf("Expected status code 200, got %d", deleteResponse.Result().StatusCode)
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
		_, err := createTestUser("123454321", "asd@ghj.io")
		if err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}

		loginResponse, err := loginTestUser("123454321")
		if err != nil {
			t.Fatalf("Failed to login test user: %v", err)
		}

		loginResponseBody, err := parseResponse[users.LoginOutputBody](loginResponse)
		if err != nil {
			t.Fatalf("Failed to decode login response: %v", err)
		}

		deleteUserHeader := makeRequestHeader(loginResponseBody.Token)
		deleteResponse := testApi.Delete(deleteUserEndpoint, deleteUserHeader)
		if deleteResponse.Result().StatusCode != http.StatusOK {
			t.Fatalf("Expected status code 200, got %d", deleteResponse.Result().StatusCode)
		}

		deleteAttempt2 := testApi.Delete(deleteUserEndpoint, deleteUserHeader)
		if deleteAttempt2.Result().StatusCode != http.StatusNotFound {
			t.Fatalf("Expected status code 404, got %d", deleteAttempt2.Result().StatusCode)
		}
	})
}
