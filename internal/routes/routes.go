package routes

import (
	"paint-api/internal/handlers"
	"paint-api/internal/handlers/brands"
	"paint-api/internal/handlers/paints"

	"github.com/danielgtaylor/huma/v2"
)

// RegisterRoutes registers the routes for the API
func RegisterRoutes(api huma.API) {
	huma.Register(api, brands.CreateOperation, brands.CreateHandler)
	huma.Register(api, brands.ListOperation, brands.ListHandler)
	huma.Register(api, brands.GetOperation, brands.GetHandler)
	huma.Register(api, brands.UpdateOperation, brands.UpdateHandler)
	huma.Register(api, brands.DeleteOperation, brands.DeleteHandler)

	huma.Register(api, paints.CreateOperation, paints.CreateHandler)
	huma.Register(api, paints.ListOperation, paints.ListHandler)
	huma.Register(api, paints.GetOperation, paints.GetHandler)
	huma.Register(api, paints.UpdateOperation, paints.UpdateHandler)
	huma.Register(api, paints.DeleteOperation, paints.DeleteHandler)

	huma.Register(api, handlers.RegisterUserOperation, handlers.RegisterUserHandler)
	huma.Register(api, handlers.LoginOperation, handlers.LoginHandler)
	huma.Register(api, handlers.RefreshTokenOperation, handlers.RefreshTokenHandler)
	huma.Register(api, handlers.ForgetUserOperation, handlers.ForgetUserHandler)

	huma.Register(api, handlers.AddToCollectionOperation, handlers.AddToCollectionHandler)
	huma.Register(api, handlers.DeleteCollectionEntryOperation, handlers.DeleteCollectionEntryHandler)
	huma.Register(api, handlers.ListPaintCollectionOperation, handlers.ListPaintCollectionHandler)
	huma.Register(api, handlers.UpdateCollectionEntryOperation, handlers.UpdateCollectionEntryHandler)
}
