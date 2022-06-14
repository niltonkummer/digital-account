package db

import (
	"digital-account/application/models"

	"gorm.io/gorm"
)

func Setup(db *gorm.DB) error {
	return db.AutoMigrate(&models.Account{}, &models.User{}, &models.TransferLock{}, &models.Transfer{})

}
