package accounts

import (
	"bytes"
	"context"
	"digital-account/application/models"
	"digital-account/application/repository"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/shopspring/decimal"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type handlerMock struct {
	repository repository.Repository
}

func (h handlerMock) Repository() repository.Repository {
	return h.repository
}

type RepositoryMock struct {
	AccountMock  func() repository.Account
	UserMock     func() repository.User
	TransferMock func() repository.Transfer
}

func (n RepositoryMock) Account() repository.Account {
	return n.AccountMock()
}

func (n RepositoryMock) User() repository.User {
	return n.UserMock()
}

func (n RepositoryMock) Transfer() repository.Transfer {
	return n.TransferMock()
}

type userRepoMock struct {
	createMock func(ctx context.Context, name, cpf, secret string) (*models.User, error)
}

func (u userRepoMock) Create(ctx context.Context, name, cpf, secret string) (*models.User, error) {
	return u.createMock(ctx, name, cpf, secret)
}

func (u userRepoMock) ByCPF(ctx context.Context, cpf string) (*models.User, error) {
	panic(nil)
}

func (u userRepoMock) Get(ctx context.Context, id int64) (*models.User, error) {
	panic(nil)
}

type accoRepoMock struct {
	getMock    func(ctx context.Context, id int64) (*models.Account, error)
	createMock func(ctx context.Context, account *models.Account) error
	listMock   func(ctx context.Context) (models.Accounts, error)
}

func (a accoRepoMock) Get(ctx context.Context, id int64) (*models.Account, error) {
	return a.getMock(ctx, id)
}

func (a accoRepoMock) Create(ctx context.Context, account *models.Account) error {
	return a.createMock(ctx, account)
}

func (a accoRepoMock) List(ctx context.Context) (models.Accounts, error) {
	return a.listMock(ctx)
}

func TestCreate(t *testing.T) {
	tt := []struct {
		name                string
		accoCreateRepoErr   error
		userCreateRepoModel *models.User
		userCreateRepoErr   error
		payload             io.Reader
		expectedStatus      int
	}{
		{name: "empty payload", expectedStatus: http.StatusBadRequest},
		{name: "password error", payload: bytes.NewBuffer([]byte(`{"name":"Fulano","secret":"","cpf":"12345678910"}`)), expectedStatus: http.StatusBadRequest},
		{name: "name error", payload: bytes.NewBuffer([]byte(`{"name":"","secret":"12345","cpf":"12345678910"}`)), expectedStatus: http.StatusBadRequest},
		{name: "cpf error", payload: bytes.NewBuffer([]byte(`{"name":"Fulano","secret":"12345","cpf":""}`)), expectedStatus: http.StatusBadRequest},
		{name: "create user error", userCreateRepoErr: errors.New("error"), payload: bytes.NewBuffer([]byte(`{"name":"Fulano","secret":"12345","cpf":"12345678910"}`)), expectedStatus: http.StatusInternalServerError},
		{name: "create account error", userCreateRepoModel: &models.User{
			Model: models.Model{
				ID: 1,
			},
			Name: "Fulano",
			CPF:  "12345678910",
		}, accoCreateRepoErr: errors.New("error"), payload: bytes.NewBuffer([]byte(`{"name":"Fulano","secret":"12345","cpf":"12345678910"}`)), expectedStatus: http.StatusInternalServerError},
		{name: "ok", userCreateRepoModel: &models.User{
			Model: models.Model{
				ID: 1,
			},
			Name: "Fulano",
			CPF:  "12345678910",
		}, payload: bytes.NewBuffer([]byte(`{"name":"Fulano","secret":"12345","cpf":"12345678910"}`)), expectedStatus: http.StatusCreated},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			a := Config(&handlerMock{
				repository: RepositoryMock{
					UserMock: func() repository.User {
						return &userRepoMock{
							createMock: func(ctx context.Context, name, cpf, secret string) (*models.User, error) {
								return tc.userCreateRepoModel, tc.userCreateRepoErr
							},
						}
					},
					AccountMock: func() repository.Account {
						return &accoRepoMock{
							createMock: func(ctx context.Context, acc *models.Account) error {
								return tc.accoCreateRepoErr
							},
						}
					},
				},
			})

			router := gin.New()
			router.POST("/", a.CreateHandler)

			req, _ := http.NewRequest("POST", fmt.Sprintf("/"), tc.payload)
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			if resp.Result().StatusCode != tc.expectedStatus {
				t.Fatalf("expected %v, got %v", tc.expectedStatus, resp.Result().StatusCode)
			}

		})
	}
}

func TestBalance(t *testing.T) {
	tt := []struct {
		name           string
		account        *models.Account
		idParam        string
		err            error
		expectedStatus int
	}{
		{"error", nil, "0", errors.New("error"), http.StatusInternalServerError},
		{"invalid param", nil, "invalid", errors.New("invalid param"), http.StatusBadRequest},
		{"ok", &models.Account{Balance: decimal.New(10, 0)}, "1", nil, http.StatusOK},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			a := Config(&handlerMock{
				repository: RepositoryMock{
					AccountMock: func() repository.Account {
						return &accoRepoMock{
							getMock: func(ctx context.Context, id int64) (*models.Account, error) {
								return tc.account, tc.err
							},
						}
					},
				},
			})

			router := gin.New()
			router.GET("/:account_id/balance", a.BalanceHandler)

			req, _ := http.NewRequest("GET", fmt.Sprintf("/%s/balance", tc.idParam), nil)
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			if resp.Result().StatusCode != tc.expectedStatus {
				t.Fatalf("expected %v, got %v", tc.expectedStatus, resp.Result().StatusCode)
			}

			if resp.Result().StatusCode == http.StatusOK {
				var resList *models.Account
				err := json.Unmarshal(resp.Body.Bytes(), &resList)
				if err != nil {
					t.Fatalf("could not unmarshal json: %v", err)
				}
				assert.EqualValues(t, tc.account, resList)
			}
		})
	}
}

func TestList(t *testing.T) {

	tt := []struct {
		name           string
		list           models.Accounts
		err            error
		expectedStatus int
	}{
		{"error", nil, errors.New("error"), http.StatusInternalServerError},
		{"ok", models.Accounts{
			&models.Account{
				Model: models.Model{
					ID:        1,
					CreatedAt: time.Time{},
					UpdatedAt: time.Time{},
					DeletedAt: gorm.DeletedAt{},
				},
				Balance: decimal.New(10, 0),
				Type:    0,
			},
			&models.Account{
				Model: models.Model{
					ID:        2,
					CreatedAt: time.Time{},
					UpdatedAt: time.Time{},
					DeletedAt: gorm.DeletedAt{},
				},
				Balance: decimal.New(20, 0),
				Type:    0,
			},
			&models.Account{
				Model: models.Model{
					ID:        3,
					CreatedAt: time.Time{},
					UpdatedAt: time.Time{},
					DeletedAt: gorm.DeletedAt{},
				},
				Balance: decimal.New(30, 0),
				Type:    0,
			},
		}, nil, http.StatusOK},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {

			a := Config(&handlerMock{
				repository: RepositoryMock{
					AccountMock: func() repository.Account {
						return &accoRepoMock{
							listMock: func(ctx context.Context) (models.Accounts, error) {
								return tc.list, tc.err
							},
						}
					},
				},
			})

			router := gin.New()
			router.GET("/list", a.ListHandler)

			req, _ := http.NewRequest("GET", "/list", nil)
			resp := httptest.NewRecorder()
			router.ServeHTTP(resp, req)

			if resp.Result().StatusCode != tc.expectedStatus {
				t.Fatalf("expected %v, got %v", tc.expectedStatus, resp.Result().StatusCode)
			}

			if resp.Result().StatusCode == http.StatusOK {
				var resList models.Accounts
				err := json.Unmarshal(resp.Body.Bytes(), &resList)
				if err != nil {
					t.Fatalf("could not unmarshal json: %v", err)
				}
				assert.EqualValues(t, tc.list, resList)
			}
		})
	}

}
