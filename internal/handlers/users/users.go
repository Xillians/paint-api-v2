package users

import (
	"github.com/danielgtaylor/huma/v2"
)

func RegisterRoutes(api huma.API) {
	huma.Register(api, forgetOperation, forgetHandler)
	huma.Register(api, refreshTokenOperation, refreshTokenHandler)
	huma.Register(api, loginOperation, loginHandler)
	huma.Register(api, registerOperation, registerHandler)
}
