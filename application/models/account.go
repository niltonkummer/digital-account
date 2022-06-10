package models

import "github.com/shopspring/decimal"

type Account struct {
	Model

	Balance decimal.Decimal `json:"balance"`
	Secret  string          `json:"-"`
	UserID  int64           `json:"-"`
	User    User            `json:"user"`
}

type Accounts []*Account
