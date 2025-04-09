package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"paint-api/internal/middleware"
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/humatest"
	"gorm.io/gorm"
)

func TestUseDb(t *testing.T) {
	// Create a mock *gorm.DB instance
	mockDb := &gorm.DB{}

	req, _ := http.NewRequest(http.MethodGet, "/demo", nil)
	w := httptest.NewRecorder()
	ctx := humatest.NewContext(&huma.Operation{}, req, w)

	middleware.UseDb(mockDb)(ctx, func(ctx huma.Context) {
		// Retrieve the *gorm.DB instance from the context
		db, ok := ctx.Context().Value(middleware.DbKey).(*gorm.DB)
		if !ok {
			t.Fatal("expected *gorm.DB to be in context")
		}
		if db != mockDb {
			t.Fatal("expected *gorm.DB to be the same as the mock instance")
		}
	})
}
