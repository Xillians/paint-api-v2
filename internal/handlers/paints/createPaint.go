package paints

import (
	"context"
	"log/slog"
	"net/http"
	"paint-api/internal/db"
	"paint-api/internal/middleware"
	"regexp"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

type CreatePaintInput struct {
	Body db.CreatePaintInput
}

type createPaintOutput struct {
	Body db.Paints
}

var createOperation = huma.Operation{
	Method: http.MethodPost,
	Path:   "/paints",
	Tags:   []string{"paints"},
}

func CreateHandler(ctx context.Context, input *CreatePaintInput) (*createPaintOutput, error) {
	if !ValidateColorCode(input.Body.ColorCode) {
		return nil, huma.NewError(http.StatusBadRequest, "Invalid color code")
	}

	connection, ok := ctx.Value(middleware.DbKey).(*gorm.DB)
	if !ok {
		slog.Error("Could not retrieve db from context")
		return nil, huma.NewError(http.StatusInternalServerError, "Failed to create paint")
	}

	paint, err := db.Paints{}.CreatePaint(connection, &input.Body)
	if err != nil {
		slog.Error("Failed to create paint", "error", err, "paint", input.Body)
		return nil, huma.NewError(http.StatusInternalServerError, "Failed to create paint")
	}
	return &createPaintOutput{Body: *paint}, nil
}

func ValidateColorCode(colorCode string) bool {
	regex := regexp.MustCompile(`^#([A-Fa-f0-9]{6})$`)
	return regex.MatchString(colorCode)
}
