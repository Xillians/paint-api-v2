package brands_test

import (
	"net/http"
	"paint-api/internal/db"
	"testing"
)

func TestDeleteBrand(t *testing.T) {
	token, err := getUserToken(testData.User.GoogleUserId)
	if err != nil {
		t.Fatalf("Failed to get user token: %v", err)
	}

	t.Run("Delete brand", func(t *testing.T) {
		testBrandName := "Test Brand"
		createResponse := createTestBrand(testBrandName, token)
		if createResponse.Result().StatusCode != http.StatusOK {
			t.Fatalf("Expected status code 201, got %d", createResponse.Result().StatusCode)
		}

		body, err := parseResponse[db.PaintBrands](createResponse)
		if err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}

		deleteResponse := deleteTestBrand(body.ID, token)
		if deleteResponse.Result().StatusCode != http.StatusOK {
			t.Fatalf("Expected status code 200, got %d", deleteResponse.Result().StatusCode)
		}
	})
}
