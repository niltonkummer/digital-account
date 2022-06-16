package accounts

import (
	"digital-account/application/api/handler"
)

type Account struct {
	handler.Handler
}

func Config(h handler.Handler) *Account {
	return &Account{
		Handler: h,
	}
}
