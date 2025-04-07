package users_test

import (
	"net/http"
	"paint-api/internal/handlers/users"
	"testing"
)

func TestRefreshToken(t *testing.T) {
	t.Run("Try to refresh token with valid token", func(t *testing.T) {
		loginResponse := loginTestUser("123456")
		if loginResponse.Result().StatusCode != http.StatusOK {
			t.Fatalf("Expected status code 200, got %d", loginResponse.Result().StatusCode)
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
		createUserResponse := createTestUser("122", "asd@asd.io")
		if createUserResponse.Result().StatusCode != http.StatusOK {
			t.Fatalf("Expected status code 200, got %d", createUserResponse.Result().StatusCode)
		}

		loginResponse := loginTestUser("122")
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

		bearer := makeRequestHeader(loginResponseBody.Token)
		refreshResponse := testApi.Get("/refresh", bearer)
		if refreshResponse.Result().StatusCode != http.StatusNotFound {
			t.Fatalf("Expected status code 401, got %d", refreshResponse.Result().StatusCode)
		}
	})
}
