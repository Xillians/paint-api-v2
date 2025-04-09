package users_test

import (
	"net/http"
	"paint-api/internal/handlers/users"
	"testing"
)

func TestForgetUser(t *testing.T) {
	deleteUserEndpoint := "/forget"
	userId := "81549300"
	t.Run("Create and delete user", func(t *testing.T) {
		createResponse := createTestUser(userId, "asd@ghj.io")
		if createResponse.Result().StatusCode != http.StatusOK {
			t.Fatalf("Expected status code 200, got %d", createResponse.Result().StatusCode)
		}

		loginResponse := loginTestUser(userId)
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
		createUserResponse := createTestUser(userId, "asd@ghj.io")
		if createUserResponse.Result().StatusCode != http.StatusOK {
			t.Fatalf("Expected status code 200, got %d", createUserResponse.Result().StatusCode)
		}

		loginResponse := loginTestUser(userId)
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

func TestRegisterUser(t *testing.T) {
	userId := "4201"
	t.Run("Successfully create a user", func(t *testing.T) {
		response := createTestUser(userId, "asd@asd.io")
		if response.Result().StatusCode != http.StatusOK {
			t.Fatalf("Expected status code 200, got %d", response.Result().StatusCode)
		}

		if response.Result().StatusCode != http.StatusOK {
			t.Fatalf("Expected status code 200, got %d", response.Result().StatusCode)
		}

		loginResponse := loginTestUser(userId)
		if loginResponse.Result().StatusCode != http.StatusOK {
			t.Fatalf("Expected status code 200, got %d", loginResponse.Result().StatusCode)
		}

		loginResponseBody, err := parseResponse[users.LoginOutputBody](loginResponse)
		if err != nil {
			t.Fatalf("Failed to decode login response: %v", err)
		}

		deleteResponse := deleteTestUser(loginResponseBody.Token)
		if deleteResponse.Result().StatusCode != http.StatusOK {
			t.Fatalf("Expected status code 200, got %d", deleteResponse.Result().StatusCode)
		}

	})
	t.Run("Attempt to create duplicate user", func(t *testing.T) {
		response := createTestUser(userId, "asd@asd.io")
		if response.Result().StatusCode != http.StatusOK {
			t.Fatalf("Expected status code 200, got %d", response.Result().StatusCode)
		}

		if response.Result().StatusCode != http.StatusOK {
			t.Fatalf("Expected status code 200, got %d", response.Result().StatusCode)
		}

		duplicateResponse := createTestUser(userId, "asd@asd.io")
		if duplicateResponse.Result().StatusCode != http.StatusConflict {
			t.Fatalf("Expected status code 409, got %d", duplicateResponse.Result().StatusCode)
		}

		loginResponse := loginTestUser(userId)
		if loginResponse.Result().StatusCode != http.StatusOK {
			t.Fatalf("Expected status code 200, got %d", loginResponse.Result().StatusCode)
		}
		loginResponseBody, err := parseResponse[users.LoginOutputBody](loginResponse)
		if err != nil {
			t.Fatalf("Failed to decode login response: %v", err)
		}

		deleteResponse := deleteTestUser(loginResponseBody.Token)
		if deleteResponse.Result().StatusCode != http.StatusOK {
			t.Fatalf("Expected status code 200, got %d", deleteResponse.Result().StatusCode)
		}
	})
	t.Run("Attempt to create user with invalid email", func(t *testing.T) {
		registerUserInput := map[string]interface{}{
			"email":   "invalid-email",
			"user_id": userId,
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
