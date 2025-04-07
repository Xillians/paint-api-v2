package brands_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"paint-api/internal/db"
	"paint-api/internal/handlers/brands"
	"testing"
)

func TestDeleteBrand(t *testing.T) {
	token, err := getUserToken(testData.User.GoogleUserId)
	if err != nil {
		t.Fatalf("Failed to get user token: %v", err)
	}

	t.Run("Delete brand", func(t *testing.T) {
		bearerToken := fmt.Sprintf("Authorization: Bearer %s", token)
		createInput := db.CreateBrandInput{
			Name: "Test Brand",
		}
		createResponse := testApi.Post("/paint-brands", bearerToken, createInput)
		if createResponse.Result().StatusCode != http.StatusOK {
			t.Fatalf("Expected status code 200, got %d", createResponse.Result().StatusCode)
		}
		var body brands.CreateBrandOutput
		err := json.NewDecoder(createResponse.Result().Body).Decode(&body.Body)
		if err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}

		deletePath := fmt.Sprintf("/paint-brands/%d", body.Body.ID)
		response := testApi.Delete(deletePath, bearerToken)
		if response.Result().StatusCode != http.StatusOK {
			t.Fatalf("Expected status code 200, got %d", response.Result().StatusCode)
		}
	})
}
