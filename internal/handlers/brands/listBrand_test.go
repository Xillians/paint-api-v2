package brands_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"paint-api/internal/handlers/brands"
	"testing"
)

func TestListBrands(t *testing.T) {
	token, err := getUserToken(testData.User.GoogleUserId)
	if err != nil {
		t.Fatalf("Failed to get user token: %v", err)
	}

	t.Run("List brands", func(t *testing.T) {
		bearerToken := fmt.Sprintf("Authorization: Bearer %s", token)
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
