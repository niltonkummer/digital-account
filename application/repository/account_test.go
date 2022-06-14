package repository

import (
	"context"
	"digital-account/application/db"
	"digital-account/application/models"
	"testing"
	"time"

	"github.com/shopspring/decimal"

	"gorm.io/driver/sqlite"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var repoTest = func() Account {

	d, _ := gorm.Open(sqlite.Open("file:memdb1?mode=memory&cache=shared"))
	_ = db.Setup(d)

	return CreateAccountRepository(d)
}

var resultGetAccount = &models.Account{
	Model: models.Model{
		ID:        1,
		CreatedAt: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
		UpdatedAt: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
	},
	Balance: decimal.NewFromFloat(0),
	Type:    0,
	Secret:  "12345",
	UserID:  1,
	User:    nil,
}

var resultListAccounts = []*models.Account{
	{
		Model: models.Model{
			ID:        2,
			CreatedAt: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		Balance: decimal.NewFromFloat(0),
		Type:    0,
		Secret:  "12345",
		UserID:  2,
		User:    nil,
	},
	{
		Model: models.Model{
			ID:        3,
			CreatedAt: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		Balance: decimal.NewFromFloat(0),
		Type:    0,
		Secret:  "12345",
		UserID:  3,
		User:    nil,
	},
	{
		Model: models.Model{
			ID:        4,
			CreatedAt: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		Balance: decimal.NewFromFloat(0),
		Type:    0,
		Secret:  "12345",
		UserID:  4,
		User:    nil,
	},
	{
		Model: models.Model{
			ID:        5,
			CreatedAt: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		Balance: decimal.NewFromFloat(0),
		Type:    0,
		Secret:  "12345",
		UserID:  5,
		User:    nil,
	},
	{
		Model: models.Model{
			ID:        6,
			CreatedAt: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		Balance: decimal.NewFromFloat(0),
		Type:    0,
		Secret:  "12345",
		UserID:  6,
		User:    nil,
	},
	{
		Model: models.Model{
			ID:        7,
			CreatedAt: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		Balance: decimal.NewFromFloat(0),
		Type:    0,
		Secret:  "12345",
		UserID:  7,
		User:    nil,
	},
	{
		Model: models.Model{
			ID:        8,
			CreatedAt: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			UpdatedAt: time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		Balance: decimal.NewFromFloat(0),
		Type:    0,
		Secret:  "12345",
		UserID:  8,
		User:    nil,
	},
}

func Test_accountRepo_Create(t *testing.T) {

	type args struct {
		ctx     context.Context
		account *models.Account
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"ok", args{
			ctx:     context.TODO(),
			account: resultGetAccount,
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if err := repoTest().Create(tt.args.ctx, tt.args.account); (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_accountRepo_Get(t *testing.T) {

	type args struct {
		ctx       context.Context
		accountID int64
	}
	tests := []struct {
		name    string
		args    args
		result  *models.Account
		wantErr bool
	}{
		{
			"ok", args{
				ctx:       context.TODO(),
				accountID: 1,
			}, resultGetAccount,
			false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			account, err := repoTest().Get(tt.args.ctx, tt.args.accountID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.Equal(t, *account, *tt.result)

		})
	}
}

func Test_accountRepo_List(t *testing.T) {

	tr := repoTest()
	for _, acc := range resultListAccounts {
		_ = tr.Create(context.TODO(), acc)
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		result  models.Accounts
		wantErr bool
	}{
		{
			"ok", args{
				ctx: context.TODO(),
			}, append([]*models.Account{resultGetAccount}, resultListAccounts...),
			false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			accounts, err := tr.List(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.Equal(t, tt.result, accounts)

		})
	}
}
