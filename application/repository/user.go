package repository

import (
	"context"
	"digital-account/application/models"

	"gorm.io/gorm"
)

type User interface {
	Create(ctx context.Context, name, cpf string, secret string) (*models.User, error)
	Auth(ctx context.Context, cpf string) (*models.User, error)
	Get(ctx context.Context, id int64) (*models.User, error)
}

type loginRepo struct {
	DB *gorm.DB
}

func (lr *loginRepo) Create(ctx context.Context, name string, cpf, secret string) (user *models.User, err error) {

	user = &models.User{
		Name:   name,
		CPF:    cpf,
		Secret: secret,
	}

	res := lr.DB.WithContext(ctx).Joins("Account").Where("cpf = ?", cpf).FirstOrCreate(user)
	err = res.Error
	return

}

func (lr *loginRepo) Auth(ctx context.Context, cpf string) (user *models.User, err error) {

	res := lr.DB.WithContext(ctx).Joins("Account").Where("cpf = ?", cpf).First(&user)
	err = res.Error
	return
}

func (lr *loginRepo) Get(ctx context.Context, id int64) (user *models.User, err error) {

	res := lr.DB.WithContext(ctx).Joins("Account").First(&user, id)
	err = res.Error
	return
}

func CreateLoginRepository(db *gorm.DB) User {

	return &loginRepo{
		DB: db,
	}
}
