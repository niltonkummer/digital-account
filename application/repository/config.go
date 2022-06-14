package repository

import (
	"gorm.io/gorm"
)

func Config(db *gorm.DB) Repository {
	return &repo{
		user:     CreateLoginRepository(db),
		account:  CreateAccountRepository(db),
		transfer: CreateTransferRepository(db),
	}
}
