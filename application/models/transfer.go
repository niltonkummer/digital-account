package models

import "github.com/shopspring/decimal"

type Transfer struct {
	Model
	AccountOrigin        Account
	AccountOriginID      int64
	AccountDestination   Account
	AccountDestinationID int64
	Amount               decimal.Decimal
}
