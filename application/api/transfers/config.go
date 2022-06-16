package transfers

import (
	"digital-account/application/api/handler"
)

type Transfers struct {
	handler.Handler
}

func Config(h handler.Handler) *Transfers {
	return &Transfers{
		Handler: h,
	}
}
