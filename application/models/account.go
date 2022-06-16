package models

import (
	"github.com/shopspring/decimal"
)

type Account struct {
	Model
	Balance decimal.Decimal `json:"balance"`
	Type    Type            `gorm:"uniqueIndex:unique_account;default:0" json:"type"`
	Secret  string          `json:"-"`
	UserID  int64           `gorm:"uniqueIndex:unique_account" json:"-"`
	User    *User           `json:"user,omitempty"`
}

type Accounts []*Account

type Type int

const (
	CurrentAccount Type = iota
	SavingsAccount
)

var typeText = [...]string{
	"current",
	"savings",
}

func (t Type) String() string {
	return typeText[t]
}
