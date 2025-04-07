package users_test

import (
	"fmt"
	"paint-api/internal/handlers/users"
	"testing"
)

func TestLogin(t *testing.T) {
	t.Run("Login with valid credentials", func(t *testing.T) {
		loginResponse, err := loginTestUser("123456")
		if err != nil {
			t.Fatalf("Failed to login test user: %v", err)
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
		if loginResponse.Result().StatusCode != 200 {
			t.Fatalf("Expected status code 200, got %d", loginResponse.Result().StatusCode)
		}
	})
	t.Run("Login with invalid credentials", func(t *testing.T) {
		loginUrl := fmt.Sprintf("/login/%s", "123")
		loginResponse := testApi.Get(loginUrl)
		if loginResponse.Result().StatusCode != 404 {
			t.Fatalf("Expected status code 404, got %d", loginResponse.Result().StatusCode)
		}
	})
}
