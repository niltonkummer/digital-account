package transfers

import (
	"github.com/shopspring/decimal"
)

type TransferRequest struct {
	//AccountOriginID      int64           `json:"account_origin_id" binding:"required"`
	AccountDestinationID int64           `json:"account_destination_id" binding:"required"`
	Amount               decimal.Decimal `json:"amount" binding:"required"`
}
