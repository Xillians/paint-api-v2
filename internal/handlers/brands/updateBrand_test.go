package brands_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"paint-api/internal/db"
	"paint-api/internal/handlers/brands"
	"testing"
)

func TestUpdateBrand(t *testing.T) {
	token, err := getUserToken(testData.User.GoogleUserId)
	if err != nil {
		t.Fatalf("Failed to get user token: %v", err)
	}

	t.Run("Update brand", func(t *testing.T) {
		bearerToken := fmt.Sprintf("Authorization: Bearer %s", token)
		payload := db.UpdateBrandInput{
			Name: "Updated Brand",
		}
		response := testApi.Put("/paint-brands/1", bearerToken, payload)
		if response.Result().StatusCode != http.StatusOK {
			t.Fatalf("Expected status code 200, got %d", response.Result().StatusCode)
		}
		var body brands.UpdateBrandOutput
		err := json.NewDecoder(response.Result().Body).Decode(&body.Body)
		if err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
	})
	t.Run("Update brand with invalid id", func(t *testing.T) {
		bearerToken := fmt.Sprintf("Authorization: Bearer %s", token)
		payload := db.UpdateBrandInput{
			Name: "Updated Brand",
		}
		response := testApi.Put("/paint-brands/100", bearerToken, payload)
		if response.Result().StatusCode != http.StatusNotFound {
			t.Fatalf("Expected status code 404, got %d", response.Result().StatusCode)
		}
	})
}
