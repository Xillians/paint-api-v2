package users_test

import (
	"net/http"
	"paint-api/internal/handlers/users"
	"testing"
)

func TestLogin(t *testing.T) {
	t.Run("Login with valid credentials", func(t *testing.T) {
		loginResponse := loginTestUser("123456")
		if loginResponse.Result().StatusCode != http.StatusOK {
			t.Fatalf("Expected status code 200, got %d", loginResponse.Result().StatusCode)
		}

		loginResponseBody, err := parseResponse[users.LoginOutputBody](loginResponse)
		if err != nil {
			t.Fatalf("Failed to decode login response: %v", err)
		}
		if loginResponseBody.Token == "" {
			t.Fatalf("Expected token in response, got empty string")
		}
		if loginResponseBody.ExpiresAt == "" {
			t.Fatalf("Expected expires_at in response, got empty string")
		}
	})
	t.Run("Login with invalid credentials", func(t *testing.T) {
		loginResponse := loginTestUser("123")
		if loginResponse.Result().StatusCode != http.StatusNotFound {
			t.Fatalf("Expected status code 404, got %d", loginResponse.Result().StatusCode)
		}
	})
}
