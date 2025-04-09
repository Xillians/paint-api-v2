package paints

import (
	"github.com/danielgtaylor/huma/v2"
)

func RegisterRoutes(api huma.API) {
	huma.Register(api, createOperation, CreateHandler)
	huma.Register(api, listOperation, ListHandler)
	huma.Register(api, getOperation, GetHandler)
	huma.Register(api, updateOperation, UpdateHandler)
	huma.Register(api, deleteOperation, DeleteHandler)
}
