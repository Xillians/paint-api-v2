package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"paint-api/internal/config"
	"paint-api/internal/db"
	"paint-api/internal/handlers/brands"
	"paint-api/internal/handlers/paint_collection"
	"paint-api/internal/handlers/paints"
	"paint-api/internal/handlers/users"
	"paint-api/internal/jwt"
	"paint-api/internal/middleware"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
)

func main() {
	apiConfig := huma.DefaultConfig("Paint API", "0.1.0")
	apiConfig.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		"bearer": {
			Type:   "http",
			Scheme: "bearer",
		},
	}
	apiConfig.Security = []map[string][]string{
		{"bearer": {}},
	}
	c := config.NewConfig()

	logLevel := c.GetLogLevel()
	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     logLevel,
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, opts))
	slog.SetDefault(logger)

	db, err := db.New(&c.DbConfig)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		panic(err)
	}

	mux := chi.NewMux()

	api := humachi.New(mux, apiConfig)
	api.UseMiddleware(middleware.UseDb(db))

	jwtService := jwt.NewJWTService(c.JwtSecret)
	api.UseMiddleware(middleware.UseJwt(*jwtService))
	api.UseMiddleware(middleware.AuthenticateRequests(api, *jwtService))

	brands.RegisterRoutes(api)
	paint_collection.RegisterRoutes(api)
	paints.RegisterRoutes(api)
	users.RegisterRoutes(api)

	slog.Info("Starting server", "port", c.HttpPort)
	err = http.ListenAndServe(fmt.Sprintf(":%d", c.HttpPort), mux)
	if err != nil {
		slog.Error("Failed to start server", "error", err)
		panic(err)
	}
}
