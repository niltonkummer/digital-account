package users

import (
	"digital-account/application/api/handler"
)

type User struct {
	handler.Handler
}

func Config(h handler.Handler) *User {
	return &User{
		Handler: h,
	}
}
