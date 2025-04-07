package users_test

import (
	"net/http"
	"paint-api/internal/handlers/users"
	"testing"
)

func TestRefreshToken(t *testing.T) {
	t.Run("Try to refresh token with valid token", func(t *testing.T) {
		loginResponse, err := loginTestUser("123456")
		if err != nil {
			t.Fatalf("Failed to login test user: %v", err)
		}

		loginResponseBody, err := parseResponse[users.LoginOutputBody](loginResponse)
		if err != nil {
			t.Fatalf("Failed to parse login response: %v", err)
		}
		bearer := makeRequestHeader(loginResponseBody.Token)

		refreshResponse := testApi.Get("/refresh", bearer)
		if refreshResponse.Result().StatusCode != http.StatusOK {
			t.Fatalf("Expected status code 200, got %d", refreshResponse.Result().StatusCode)
		}
	})
	t.Run("Refresh with invalid bearer", func(t *testing.T) {
		bearer := makeRequestHeader("asd")
		refreshResponse := testApi.Get("/refresh", bearer)
		if refreshResponse.Result().StatusCode != http.StatusUnauthorized {
			t.Fatalf("Expected status code 401, got %d", refreshResponse.Result().StatusCode)
		}
	})
	t.Run("Refresh with token where user is deleted", func(t *testing.T) {
		_, err := createTestUser("122", "asd@asd.io")
		if err != nil {
			t.Fatalf("Failed to create test user: %v", err)
		}

		loginResponse, err := loginTestUser("122")
		if err != nil {
			t.Fatalf("Failed to login test user: %v", err)
		}
		loginResponseBody, err := parseResponse[users.LoginOutputBody](loginResponse)
		if err != nil {
			t.Fatalf("Failed to decode login response: %v", err)
		}
		_, err = deleteTestUser(loginResponseBody.Token)
		if err != nil {
			t.Fatalf("Failed to delete test user: %v", err)
		}

		bearer := makeRequestHeader(loginResponseBody.Token)
		refreshResponse := testApi.Get("/refresh", bearer)
		if refreshResponse.Result().StatusCode != http.StatusNotFound {
			t.Fatalf("Expected status code 401, got %d", refreshResponse.Result().StatusCode)
		}
	})
}
