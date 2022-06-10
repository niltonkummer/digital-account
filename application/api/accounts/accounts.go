package accounts

import (
	"digital-account/application/api/common"
	"digital-account/application/models"
	"github.com/shopspring/decimal"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"

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

	user, err := common.UserFromContext(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	account := &models.Account{
		Secret:  string(bcryptPass),
		Balance: decimal.Zero,
		UserID:  user.ID,
	}

	err = a.app.Repository.Account().Create(c.Request.Context(), account)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, account)
}

// GetHandler retrieves an account by id
func (a *Account) GetHandler(c *gin.Context) {

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

	c.JSON(http.StatusOK, account)
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
