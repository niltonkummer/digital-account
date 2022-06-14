package transfers

import (
	"digital-account/application/config"
)

type Transfers struct {
	app *config.App
}

func Config(app *config.App) *Transfers {
	return &Transfers{
		app: app,
	}
}
