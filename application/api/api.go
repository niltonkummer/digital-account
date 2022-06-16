package api

import (
	"digital-account/application/api/accounts"
	"digital-account/application/api/handler"
	"digital-account/application/api/transfers"
	"digital-account/application/api/users"
	"digital-account/application/config"
)

type API struct {
	AccountsService  *accounts.Account
	UserService      *users.User
	TransfersService *transfers.Transfers
}

func SetupRoutes(app *config.App) *API {

	h := handler.New(app.DB)
	a := &API{
		AccountsService:  accounts.Config(h),
		UserService:      users.Config(h),
		TransfersService: transfers.Config(h),
	}

	a.Routes(app)

	return a
}
