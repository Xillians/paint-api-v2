package users_test

import (
	"encoding/json"
	"fmt"
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

func createTestUser(userId string, email string) *httptest.ResponseRecorder {
	registerUserInput := map[string]interface{}{
		"email":   email,
		"user_id": userId,
	}
	createResponse := testApi.Post("/register", registerUserInput)
	return createResponse
}

func loginTestUser(userId string) *httptest.ResponseRecorder {
	loginUrl := fmt.Sprintf("/login/%s", userId)
	loginResponse := testApi.Get(loginUrl)
	return loginResponse
}

func deleteTestUser(userId string) *httptest.ResponseRecorder {
	bearer := makeRequestHeader(userId)
	deleteResponse := testApi.Delete("/forget", bearer)
	return deleteResponse
}
