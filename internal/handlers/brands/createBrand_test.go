package brands_test

import (
	"net/http"
	"paint-api/internal/db"
	"testing"
)

func TestCreateBrand(t *testing.T) {
	token, err := getUserToken(testData.User.GoogleUserId)
	if err != nil {
		t.Fatalf("Failed to get user token: %v", err)
	}

	t.Run("Create brand", func(t *testing.T) {
		createBrandResponse := createTestBrand("Test Brand", token)
		if createBrandResponse.Result().StatusCode != http.StatusOK {
			t.Fatalf("Expected status code 200, got %d", createBrandResponse.Result().StatusCode)
		}

		body, err := parseResponse[db.PaintBrands](createBrandResponse)
		if err != nil {
			t.Fatalf("Failed to parse create brand response: %v", err)
		}

		if body.Name != "Test Brand" {
			t.Fatalf("Expected name to be Test Brand, got %s", body.Name)
		}
	})
}
