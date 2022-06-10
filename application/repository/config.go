package repository

import (
	"gorm.io/gorm"
)

func Config(db *gorm.DB) Repository {
	return &repo{
		login:   CreateLoginRepository(db),
		account: CreateAccountRepository(db),
	}
}
