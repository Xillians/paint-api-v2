package routes

import (
	"paint-api/internal/handlers"

	"github.com/danielgtaylor/huma/v2"
)

// RegisterRoutes registers the routes for the API
func RegisterRoutes(api huma.API) {
	huma.Register(api, handlers.CreatePaintBrandOperation, handlers.CreatePaintBrandHandler)
	huma.Register(api, handlers.ListPaintBrandsOperation, handlers.ListPaintBrandsHandler)
	huma.Register(api, handlers.GetPaintBrandOperation, handlers.GetPaintBrandHandler)
	huma.Register(api, handlers.UpdatePaintBrandOperation, handlers.UpdatePaintBrandHandler)
	huma.Register(api, handlers.DeletePaintBrandOperation, handlers.DeletePaintBrandHandler)
}
