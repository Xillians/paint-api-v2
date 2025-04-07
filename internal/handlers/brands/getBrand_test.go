package brands_test

import (
	"net/http"
	"paint-api/internal/db"
	"testing"
)

func TestGetBrand(t *testing.T) {
	token, err := getUserToken(testData.User.GoogleUserId)
	if err != nil {
		t.Fatalf("Failed to get user token: %v", err)
	}

	t.Run("Get brand", func(t *testing.T) {
		response := getTestBrand(testData.Brand.ID, token)
		if response.Result().StatusCode != http.StatusOK {
			t.Fatalf("Expected status code 200, got %d", response.Result().StatusCode)
		}
		responseBody, err := parseResponse[db.PaintBrands](response)
		if err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}
		if responseBody.ID != 1 {
			t.Fatalf("Expected brand ID 1, got %d", responseBody.ID)
		}
		if responseBody.Name != testData.Brand.Name {
			t.Fatalf("Expected brand name %s, got %s", testData.Brand.Name, responseBody.Name)
		}
	})
	t.Run("Get brand with invalid id", func(t *testing.T) {
		response := getTestBrand(100, token)
		if response.Result().StatusCode != http.StatusNotFound {
			t.Fatalf("Expected status code 404, got %d", response.Result().StatusCode)
		}
	})
}
