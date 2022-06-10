package repository

import (
	"context"
	"digital-account/application/models"
	"gorm.io/gorm"
)

type Login interface {
	Auth(ctx context.Context, cpf string, secret string) (*models.User, error)
}

type loginRepo struct {
	DB *gorm.DB
}

func (lr *loginRepo) Auth(ctx context.Context, cpf, secret string) (user *models.User, err error) {

	res := lr.DB.WithContext(ctx).Where("cpf = ?", cpf).First(&user)
	err = res.Error
	return
}

func CreateLoginRepository(db *gorm.DB) Login {

	return &loginRepo{
		DB: db,
	}
}
