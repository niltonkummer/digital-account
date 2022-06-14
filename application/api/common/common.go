package common

import (
	"digital-account/application/models"
	"errors"

	"github.com/gin-gonic/gin"
)

const IdentityKey = "jti"

func UserFromContext(c *gin.Context) (*models.User, error) {
	u, exists := c.Get(IdentityKey)
	if !exists {
		return nil, errors.New("users not exists")
	}
	user, _ := u.(*models.User)
	return user, nil
}
