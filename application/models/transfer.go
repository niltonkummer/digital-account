package models

import (
	"github.com/shopspring/decimal"
)

type Transfer struct {
	Model
	AccountOrigin        Account         `json:"-"`
	AccountOriginID      int64           `json:"account_origin_id"`
	AccountDestination   Account         `json:"-"`
	AccountDestinationID int64           `json:"account_destination_id"`
	Amount               decimal.Decimal `json:"amount"`
}

type Transfers []Transfer
