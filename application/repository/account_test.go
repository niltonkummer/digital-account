package repository

import (
	"digital-account/application/db"
	"digital-account/application/models"
	"reflect"
	"testing"
	"time"

	"gorm.io/driver/sqlite"

	"gorm.io/gorm"
)

var DB = func() *gorm.DB {
	db, _ := db.Config(sqlite.Open("../../db/test_data.sqlite"))
	return db
}

func Test_accountRepo_GetAccount(t *testing.T) {
	type fields struct {
		DB *gorm.DB
	}
	type args struct {
		id int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *models.Account
		wantErr bool
	}{
		{
			name: "found",
			fields: fields{
				DB: DB(),
			},
			args: args{
				id: 1,
			},
			want: &models.Account{
				Model: models.Model{
					ID:        1,
					CreatedAt: time.Date(2021, 4, 19, 1, 12, 35, 0, time.UTC),
					UpdatedAt: time.Date(2021, 4, 19, 1, 12, 39, 0, time.UTC),
				},
			},
			wantErr: false,
		},
		{
			name: "not_found_or_db_error",
			fields: fields{
				DB: DB(),
			},
			args: args{
				id: 12,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &accountRepo{
				DB: tt.fields.DB,
			}
			got, err := a.Get(nil, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_accountRepo_ListAccounts(t *testing.T) {
	type fields struct {
		DB *gorm.DB
	}
	tests := []struct {
		name    string
		fields  fields
		want    models.Accounts
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				DB: DB(),
			},
			want: models.Accounts{
				{
					Model: models.Model{
						ID:        1,
						CreatedAt: time.Date(2021, 4, 19, 1, 12, 35, 0, time.UTC),
						UpdatedAt: time.Date(2021, 4, 19, 1, 12, 39, 0, time.UTC),
						DeletedAt: gorm.DeletedAt{},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &accountRepo{
				DB: tt.fields.DB,
			}
			got, err := a.List(nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List() got = %v, want %v", got, tt.want)
			}
		})
	}
}
