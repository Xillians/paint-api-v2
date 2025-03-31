package routes

import (
	"paint-api/internal/handlers"
	"paint-api/internal/handlers/brands"
	"paint-api/internal/handlers/paints"
	"paint-api/internal/handlers/users"

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

	huma.Register(api, users.RegisterOperation, users.RegisterHandler)
	huma.Register(api, users.LoginOperation, users.LoginHandler)
	huma.Register(api, users.RefreshTokenOperation, users.RefreshTokenHandler)
	huma.Register(api, users.ForgetOperation, users.ForgetHandler)

	huma.Register(api, handlers.AddToCollectionOperation, handlers.AddToCollectionHandler)
	huma.Register(api, handlers.DeleteCollectionEntryOperation, handlers.DeleteCollectionEntryHandler)
	huma.Register(api, handlers.ListPaintCollectionOperation, handlers.ListPaintCollectionHandler)
	huma.Register(api, handlers.UpdateCollectionEntryOperation, handlers.UpdateCollectionEntryHandler)
}
