package paints_test

import (
	"context"
	"errors"
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
