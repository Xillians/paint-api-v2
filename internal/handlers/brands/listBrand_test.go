package brands_test

import (
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
		response := getTestBrands(token)
		if response.Result().StatusCode != http.StatusOK {
			t.Fatalf("Expected status code 200, got %d", response.Result().StatusCode)
		}
		responseBody, err := parseResponse[brands.ListBrandOutputBody](response)
		if err != nil {
			t.Fatalf("Failed to parse response: %v", err)
		}
		if len(responseBody.Brands) == 0 {
			t.Fatalf("Expected non-empty response, got empty")
		}
	})
}
