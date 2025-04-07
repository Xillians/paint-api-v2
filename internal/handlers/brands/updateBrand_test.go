package brands_test

import (
	"net/http"
	"paint-api/internal/db"
	"testing"
)

func TestUpdateBrand(t *testing.T) {
	token, err := getUserToken(testData.User.GoogleUserId)
	if err != nil {
		t.Fatalf("Failed to get user token: %v", err)
	}

	t.Run("Update brand", func(t *testing.T) {
		newBrandName := "Updated Brand"
		response := updateTestBrand(testData.Brand.ID, newBrandName, token)
		if response.Result().StatusCode != http.StatusOK {
			t.Fatalf("Expected status code 200, got %d", response.Result().StatusCode)
		}

		responseBody, err := parseResponse[db.PaintBrands](response)
		if err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}
		if responseBody.ID != testData.Brand.ID {
			t.Fatalf("Expected brand ID %d, got %d", testData.Brand.ID, responseBody.ID)
		}
		if responseBody.Name != newBrandName {
			t.Fatalf("Expected brand name %s, got %s", newBrandName, responseBody.Name)
		}
	})
	t.Run("Update brand with invalid id", func(t *testing.T) {
		newBrandName := "Updated Brand"
		response := updateTestBrand(100, newBrandName, token)
		if response.Result().StatusCode != http.StatusNotFound {
			t.Fatalf("Expected status code 404, got %d", response.Result().StatusCode)
		}
	})
}
