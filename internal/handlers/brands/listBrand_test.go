package brands_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"paint-api/internal/handlers/brands"
	"paint-api/internal/handlers/users"
	"testing"
)

func TestListBrands(t *testing.T) {
	loginPath := fmt.Sprintf("/login/%s", testData.User.GoogleUserId)
	loginResponse := testApi.Get(loginPath)
	if loginResponse.Result().StatusCode != http.StatusOK {
		t.Fatalf("Expected status code 200, got %d", loginResponse.Result().StatusCode)
	}
	var body users.LoginOutputBody
	err := json.NewDecoder(loginResponse.Result().Body).Decode(&body)
	if err != nil {
		t.Fatalf("Failed to decode login response: %v", err)
	}

	t.Run("List brands", func(t *testing.T) {
		bearerToken := fmt.Sprintf("Authorization: Bearer %s", body.Token)
		response := testApi.Get("/paint-brands", bearerToken)
		if response.Result().StatusCode != http.StatusOK {
			t.Fatalf("Expected status code 200, got %d", response.Result().StatusCode)
		}
		var body brands.ListBrandOutput
		err := json.NewDecoder(response.Result().Body).Decode(&body.Body)
		if err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
	})
}
