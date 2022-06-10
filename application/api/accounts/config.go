package accounts

import (
	"digital-account/application/config"
)

type Account struct {
	app *config.App
}

func Config(app *config.App) *Account {
	return &Account{
		app: app,
	}
}
