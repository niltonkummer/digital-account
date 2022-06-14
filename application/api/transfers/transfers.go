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
	err := c.BindJSON(&req)
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

	res, err := a.app.Repository.Transfer().Create(
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

	accounts, err := a.app.Repository.Transfer().List(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, accounts)
}
