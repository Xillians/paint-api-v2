package brands_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"paint-api/internal/db"
	"paint-api/internal/handlers/brands"
	"testing"
)

func makeRequestHeader(token string) string {
	return fmt.Sprintf("Authorization: Bearer %s", token)
}
func parseResponse[T any](response *httptest.ResponseRecorder) (*T, error) {
	var body *T
	err := json.NewDecoder(response.Result().Body).Decode(&body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

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
