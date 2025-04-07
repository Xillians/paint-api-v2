package users_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"

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

func parseResponse[T any](response *httptest.ResponseRecorder) (*T, error) {
	var body *T
	err := json.NewDecoder(response.Result().Body).Decode(&body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func makeRequestHeader(header string) string {
	return fmt.Sprintf("Authorization: Bearer %s", header)
}

func createTestUser(userId string, email string) (*httptest.ResponseRecorder, error) {
	registerUserInput := map[string]interface{}{
		"email":   email,
		"user_id": userId,
	}
	createResponse := testApi.Post("/register", registerUserInput)
	if createResponse.Result().StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected status code 200, got %d", createResponse.Result().StatusCode)
	}
	return createResponse, nil
}

func loginTestUser(userId string) (*httptest.ResponseRecorder, error) {
	loginUrl := fmt.Sprintf("/login/%s", userId)
	loginResponse := testApi.Get(loginUrl)
	if loginResponse.Result().StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected status code 200, got %d", loginResponse.Result().StatusCode)
	}
	return loginResponse, nil
}
func deleteTestUser(userId string) (*httptest.ResponseRecorder, error) {
	bearer := makeRequestHeader(userId)
	deleteResponse := testApi.Delete("/forget", bearer)
	if deleteResponse.Result().StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected status code 200, got %d", deleteResponse.Result().StatusCode)
	}
	return deleteResponse, nil
}
