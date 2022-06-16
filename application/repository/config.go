package repository

import (
	"gorm.io/gorm"
)

func Config(db *gorm.DB) Repository {
	return &repo{
		user:     CreateUserRepository(db),
		account:  CreateAccountRepository(db),
		transfer: CreateTransferRepository(db),
	}
}
