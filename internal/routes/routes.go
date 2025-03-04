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

	huma.Register(api, handlers.CreatePaintOperation, handlers.CreatePaintHandler)
	huma.Register(api, handlers.ListPaintsOperation, handlers.ListPaintsHandler)
	huma.Register(api, handlers.GetPaintsOperation, handlers.GetPaintHandler)
	huma.Register(api, handlers.UpdatePaintOperation, handlers.UpdatePaintHandler)
	huma.Register(api, handlers.DeletePaintOperation, handlers.DeletePaintHandler)

	huma.Register(api, handlers.CreateUserOperation, handlers.CreateUserHandler)
	huma.Register(api, handlers.ListUsersOperation, handlers.ListUsersHandler)
	huma.Register(api, handlers.GetUsersOperation, handlers.GetUserHandler)
	huma.Register(api, handlers.DeleteUserOperation, handlers.DeleteUserHandler)
}
