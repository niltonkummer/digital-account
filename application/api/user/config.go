package user

import (
	"digital-account/application/config"
)

type User struct {
	app *config.App
}

func Config(app *config.App) *User {
	return &User{
		app: app,
	}
}
