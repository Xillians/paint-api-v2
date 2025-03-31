package testutils

import (
	"paint-api/internal/handlers/brands"
	"paint-api/internal/handlers/paint_collection"
	"paint-api/internal/handlers/paints"
	"paint-api/internal/handlers/users"
	"paint-api/internal/jwt"
	"paint-api/internal/middleware"
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/humatest"
)

func MakeTestApi(t *testing.T) (humatest.TestAPI, *TestData, func()) {
	apiConfig := huma.DefaultConfig("Paint API", "0.1.0")

	// Open the test database connection
	connection, dbCleanup := OpenTestConnection()

	// Create a cleanup function that includes database cleanup
	apiCleanup := func() {
		// Call the database cleanup function
		dbCleanup()
	}

	_, testApi := humatest.New(t, apiConfig)
	testApi.UseMiddleware(middleware.UseDb(connection))

	// Initialize JWT service
	jwtService := jwt.NewJWTService("some_secret")
	testApi.UseMiddleware(middleware.UseJwt(*jwtService))
	testApi.UseMiddleware(middleware.AuthenticateRequests(testApi, *jwtService))

	// Register routes
	brands.RegisterRoutes(testApi)
	paints.RegisterRoutes(testApi)
	users.RegisterRoutes(testApi)
	paint_collection.RegisterRoutes(testApi)

	testData, err := MakeTestData(connection)
	if err != nil {
		t.Fatalf("Failed to create test data: %v", err)
	}

	return testApi, testData, apiCleanup
}
