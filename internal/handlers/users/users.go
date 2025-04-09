package users

import (
	"github.com/danielgtaylor/huma/v2"
)

func RegisterRoutes(api huma.API) {
	huma.Register(api, forgetOperation, ForgetHandler)
	huma.Register(api, refreshTokenOperation, RefreshTokenHandler)
	huma.Register(api, loginOperation, LoginHandler)
	huma.Register(api, registerOperation, RegisterHandler)
}
