package accounts

import (
	"digital-account/application/models"
	"net/http"
	"strconv"

	"github.com/shopspring/decimal"
	"golang.org/x/crypto/bcrypt"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
)

// CreateHandler saves a new account
func (a *Account) CreateHandler(c *gin.Context) {

	req := AccountRequest{}
	err := c.BindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	bcryptPass, err := bcrypt.GenerateFromPassword([]byte(req.Secret), bcrypt.DefaultCost)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	user, err := a.app.Repository.User().Create(c, req.Name, req.CPF, string(bcryptPass))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	account := &models.Account{
		Balance: decimal.Zero,
		Secret:  string(bcryptPass),
		UserID:  user.ID,
		User:    user,
	}

	err = a.app.Repository.Account().Create(c.Request.Context(), account)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, account)
}

// BalanceHandler retrieves the balance from a account
func (a *Account) BalanceHandler(c *gin.Context) {

	idStr := c.Param("account_id")

	if !govalidator.IsNumeric(idStr) {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid param"})
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	account, errAcc := a.app.Repository.Account().Get(c.Request.Context(), id)
	if errAcc != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": errAcc})
		return
	}

	c.JSON(http.StatusOK, gin.H{"balance": account.Balance})
}

// ListHandler retrieves a list of accounts
func (a *Account) ListHandler(c *gin.Context) {

	accounts, err := a.app.Repository.Account().List(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, accounts)
}
