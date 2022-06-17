package transfers

import (
	"digital-account/application/api/common"
	"digital-account/application/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateHandler creates a new internal transfer
func (a *Transfers) CreateHandler(c *gin.Context) {

	req := TransferRequest{}
	err := c.MustBindWith(&req, &req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := common.UserFromContext(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	if user.Account.ID == req.AccountDestinationID {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "transfer to the same account is not allowed"})
		return
	}

	//hundred := decimal.New(100, 0)
	// req.Amount = decimal.NewFromInt(req.Amount.Mul(hundred).IntPart()).Div(hundred)

	res, err := a.Repository().Transfer().Create(
		c,
		&models.Transfer{
			AccountOriginID:      user.Account.ID,
			AccountDestinationID: req.AccountDestinationID,
			Amount:               req.Amount,
		},
	)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusCreated, res)
}

// ListHandler retrieves a list of transfers
func (a *Transfers) ListHandler(c *gin.Context) {

	user, err := common.UserFromContext(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	accounts, err := a.Repository().Transfer().List(c, user.ID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, accounts)
}
