package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"paint-api/internal/jwt"
	"paint-api/internal/middleware"
	"paint-api/internal/testutils"
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/humatest"
)

func TestUseJwt(t *testing.T) {
	jwtService := jwt.NewJWTService("some_secret")

	t.Run("add UseJwt", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/demo", nil)
		w := httptest.NewRecorder()
		ctx := humatest.NewContext(&huma.Operation{}, req, w)

		middleware.UseJwt(*jwtService)(ctx, func(ctx huma.Context) {
			// Retrieve the *gorm.DB instance from the context
			service, ok := ctx.Context().Value(middleware.JwtKey).(jwt.JWTService)
			if !ok {
				t.Fatal("expected *gorm.DB to be in context")
			}
			if service != *jwtService {
				t.Fatal("expected *gorm.DB to be the same as the mock instance")
			}
		})
	})
}
func TestAuthenticateRequest(t *testing.T) {
	jwtService := jwt.NewJWTService("some_secret")
	api, _, _ := testutils.MakeTestApi(t)

	t.Run("successful request", func(t *testing.T) {
		token, err := jwtService.GenerateToken("123", "administrator")
		if err != nil {
			t.Fatalf("failed to generate token: %v", err)
		}

		req, _ := http.NewRequest(http.MethodGet, "/demo", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		ctx := humatest.NewContext(&huma.Operation{}, req, w)

		middleware.AuthenticateRequests(api, *jwtService)(ctx, func(ctx huma.Context) {

		})
	})
	t.Run("/login request", func(t *testing.T) {
		token, err := jwtService.GenerateToken("123", "administrator")
		if err != nil {
			t.Fatalf("failed to generate token: %v", err)
		}

		req, _ := http.NewRequest(http.MethodGet, "/login", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		ctx := humatest.NewContext(&huma.Operation{}, req, w)

		middleware.AuthenticateRequests(api, *jwtService)(ctx, func(ctx huma.Context) {

		})
	})
	t.Run("register request", func(t *testing.T) {
		token, err := jwtService.GenerateToken("123", "administrator")
		if err != nil {
			t.Fatalf("failed to generate token: %v", err)
		}

		req, _ := http.NewRequest(http.MethodGet, "/register", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		ctx := humatest.NewContext(&huma.Operation{}, req, w)

		middleware.AuthenticateRequests(api, *jwtService)(ctx, func(ctx huma.Context) {

		})
	})
	t.Run("missing token", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/login", nil)
		req.Header.Set("Authorization", "Bearer abc")
		w := httptest.NewRecorder()

		ctx := humatest.NewContext(&huma.Operation{}, req, w)

		middleware.AuthenticateRequests(api, *jwtService)(ctx, func(ctx huma.Context) {

		})
	})
	t.Run("missing user_id", func(t *testing.T) {
		token, err := jwtService.GenerateToken("", "administrator")
		if err != nil {
			t.Fatalf("failed to generate token: %v", err)
		}

		req, _ := http.NewRequest(http.MethodGet, "/register", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		ctx := humatest.NewContext(&huma.Operation{}, req, w)

		middleware.AuthenticateRequests(api, *jwtService)(ctx, func(ctx huma.Context) {

		})
	})
	t.Run("missing role", func(t *testing.T) {
		token, err := jwtService.GenerateToken("123", "")
		if err != nil {
			t.Fatalf("failed to generate token: %v", err)
		}

		req, _ := http.NewRequest(http.MethodGet, "/register", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		ctx := humatest.NewContext(&huma.Operation{}, req, w)

		middleware.AuthenticateRequests(api, *jwtService)(ctx, func(ctx huma.Context) {

		})
	})
}
