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

	huma.Register(api, handlers.RegisterUserOperation, handlers.RegisterUserHandler)
	huma.Register(api, handlers.LoginOperation, handlers.LoginHandler)
	huma.Register(api, handlers.RefreshTokenOperation, handlers.RefreshTokenHandler)
	huma.Register(api, handlers.ForgetUserOperation, handlers.ForgetUserHandler)

	huma.Register(api, handlers.AddToCollectionOperation, handlers.AddToCollectionHandler)
	huma.Register(api, handlers.DeleteCollectionEntryOperation, handlers.DeleteCollectionEntryHandler)
	huma.Register(api, handlers.ListPaintCollectionOperation, handlers.ListPaintCollectionHandler)
	huma.Register(api, handlers.UpdateCollectionEntryOperation, handlers.UpdateCollectionEntryHandler)
}
