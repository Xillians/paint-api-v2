package brands_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"paint-api/internal/db"
	"paint-api/internal/handlers/brands"
	"paint-api/internal/handlers/users"
	"testing"
)

func TestCreateBrand(t *testing.T) {
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

	t.Run("Create brand", func(t *testing.T) {
		bearerToken := fmt.Sprintf("Authorization: Bearer %s", body.Token)
		payload := db.CreateBrandInput{
			Name: "Test Brand",
		}
		response := testApi.Post("/paint-brands", bearerToken, payload)
		if response.Result().StatusCode != http.StatusOK {
			t.Fatalf("Expected status code 200, got %d", response.Result().StatusCode)
		}
		var body brands.CreateBrandOutput
		err := json.NewDecoder(response.Result().Body).Decode(&body.Body)
		if err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		if body.Body.Name != payload.Name {
			t.Fatalf("Expected name to be %s, got %s", payload.Name, body.Body.Name)
		}
	})
}
