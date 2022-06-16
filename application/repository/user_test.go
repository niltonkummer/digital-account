package repository

import (
	"context"
	"digital-account/application/db"
	"digital-account/application/models"
	"fmt"
	"testing"
	"time"

	"gorm.io/driver/sqlite"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var UserRepoTest = func() User {

	time.Local = time.UTC
	d, _ := gorm.Open(sqlite.Open("file:memdb1?mode=memory&cache=shared"))
	_ = db.Setup(d.Debug())

	return CreateUserRepository(d)
}

func TestCreate(t *testing.T) {

	ur := UserRepoTest()

	type args struct {
		ctx    context.Context
		name   string
		cpf    string
		secret string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"ok", args{
			name:   "Fulano",
			cpf:    "99911133312",
			secret: "12345",
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			createdUser, err := ur.Create(tt.args.ctx, tt.args.name, tt.args.cpf, tt.args.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
			}

			gotUser, err := ur.ByCPF(tt.args.ctx, tt.args.cpf)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
			}
			gotUser.Account = nil

			assert.EqualValues(t, createdUser, gotUser, "Get(%v, %v)", tt.args.ctx, tt.args.cpf)

			gotUserByID, err := ur.Get(tt.args.ctx, gotUser.ID)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
			}
			gotUserByID.Account = nil
			assert.EqualValues(t, gotUserByID, gotUser, "Get(%v, %v)", tt.args.ctx, tt.args.cpf)

		})
	}
}

func TestGet(t *testing.T) {
	type fields struct {
		DB *gorm.DB
	}
	type args struct {
		ctx context.Context
		id  int64
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantUser *models.User
		wantErr  assert.ErrorAssertionFunc
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			gotUser, err := UserRepoTest().Get(tt.args.ctx, tt.args.id)
			if !tt.wantErr(t, err, fmt.Sprintf("Get(%v, %v)", tt.args.ctx, tt.args.id)) {
				return
			}
			assert.Equalf(t, tt.wantUser, gotUser, "Get(%v, %v)", tt.args.ctx, tt.args.id)
		})
	}
}
