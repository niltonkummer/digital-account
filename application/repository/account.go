package repository

import (
	"context"
	"digital-account/application/models"

	"gorm.io/gorm"
)

type Account interface {
	Get(ctx context.Context, id int64) (*models.Account, error)
	Create(ctx context.Context, account *models.Account) error
	List(ctx context.Context) (models.Accounts, error)
}

type accountRepo struct {
	DB *gorm.DB
}

func (a *accountRepo) Create(ctx context.Context, account *models.Account) error {

	res := a.DB.WithContext(ctx).Create(account)
	return res.Error
}

func (a *accountRepo) List(ctx context.Context) (accounts models.Accounts, err error) {

	res := a.DB.WithContext(ctx).Scopes(Paginate(ctx)).Joins("User").Find(&accounts)
	err = res.Error

	return
}

func (a *accountRepo) Get(ctx context.Context, id int64) (account *models.Account, err error) {

	tx := a.DB.WithContext(ctx).Joins("User").First(&account, id)
	err = tx.Error

	return
}

func CreateAccountRepository(db *gorm.DB) Account {
	return &accountRepo{
		DB: db,
	}
}
