package brands_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"

	"paint-api/internal/db"
	"paint-api/internal/handlers/users"
	"paint-api/internal/middleware"
	"paint-api/internal/testutils"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
)

var testApi humatest.TestAPI
var cleanUp func()
var testData *testutils.TestData

func TestMain(m *testing.M) {
	api, data, apiCleanup := testutils.MakeTestApi(&testing.T{})
	testApi = api
	testData = data

	cleanUp = func() {
		apiCleanup()
	}

	code := m.Run()

	cleanUp()
	os.Exit(code)
}

// Integration test helpers

func createTestBrand(brandName string, token string) *httptest.ResponseRecorder {
	brandInput := db.CreateBrandInput{
		Name: brandName,
	}
	response := testApi.Post("/paint-brands", brandInput, makeRequestHeader(token))
	return response
}
func updateTestBrand(brandID int, brandName string, token string) *httptest.ResponseRecorder {
	brandInput := db.UpdateBrandInput{
		Name: brandName,
	}
	path := fmt.Sprintf("/paint-brands/%d", brandID)
	response := testApi.Put(path, brandInput, makeRequestHeader(token))
	return response
}
func deleteTestBrand(brandID int, token string) *httptest.ResponseRecorder {
	path := fmt.Sprintf("/paint-brands/%d", brandID)
	response := testApi.Delete(path, makeRequestHeader(token))
	return response
}
func getTestBrand(brandID int, token string) *httptest.ResponseRecorder {
	path := fmt.Sprintf("/paint-brands/%d", brandID)
	response := testApi.Get(path, makeRequestHeader(token))
	return response
}
func getTestBrands(token string) *httptest.ResponseRecorder {
	response := testApi.Get("/paint-brands", makeRequestHeader(token))
	return response
}
func userLogin(userId string) *httptest.ResponseRecorder {
	loginPath := fmt.Sprintf("/login/%s", userId)
	loginResponse := testApi.Get(loginPath)
	return loginResponse
}
func getUserToken(userId string) (string, error) {
	loginResponse := userLogin(userId)
	if loginResponse.Result().StatusCode != http.StatusOK {
		return "", errors.New("failed to login")
	}
	body, err := parseResponse[users.LoginOutputBody](loginResponse)
	if err != nil {
		return "", fmt.Errorf("failed to parse login response: %w", err)
	}
	return body.Token, nil
}

// Unit test helpers

func createClosedDBContext() (context.Context, error) {
	ctx := context.Background()

	connection, _ := testutils.OpenTestConnection()
	sql, err := connection.DB()
	if err != nil {
		return nil, errors.New("failed to get DB from connection")
	}
	sql.Close()

	ctx = context.WithValue(ctx, middleware.DbKey, connection)
	return ctx, nil
}
