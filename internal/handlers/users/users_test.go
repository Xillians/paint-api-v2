package users_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http/httptest"
	"os"

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

// Integration test helpers

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
