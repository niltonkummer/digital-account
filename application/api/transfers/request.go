package transfers

import (
	"net/http"

	"github.com/gin-gonic/gin/binding"

	"github.com/shopspring/decimal"
)

type TransferRequest struct {
	AccountDestinationID int64           `json:"account_destination_id" binding:"required"`
	Amount               decimal.Decimal `json:"amount" binding:"required"`
}

func (t TransferRequest) Name() string {
	return "TransferRequest"
}

func (t *TransferRequest) Bind(request *http.Request, i interface{}) error {

	err := binding.JSON.Bind(request, i)
	if err != nil {
		return err
	}

	hundred := decimal.New(100, 0)
	t.Amount = decimal.NewFromInt(t.Amount.Mul(hundred).IntPart()).Div(hundred)

	return nil
}
