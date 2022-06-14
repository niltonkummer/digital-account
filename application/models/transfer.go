package models

import "github.com/shopspring/decimal"

type Transfer struct {
	Model
	AccountOrigin        Account         `json:"account_origin,omitempty"`
	AccountOriginID      int64           `json:"-"`
	AccountDestination   Account         `json:"account_destination,omitempty"`
	AccountDestinationID int64           `json:"-"`
	Amount               decimal.Decimal `json:"amount"`
}

type Transfers []Transfer
