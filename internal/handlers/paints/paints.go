package paints

import (
	"github.com/danielgtaylor/huma/v2"
)

func RegisterRoutes(api huma.API) {
	huma.Register(api, createOperation, createHandler)
	huma.Register(api, listOperation, listHandler)
	huma.Register(api, getOperation, getHandler)
	huma.Register(api, updateOperation, updateHandler)
	huma.Register(api, deleteOperation, deleteHandler)
}
