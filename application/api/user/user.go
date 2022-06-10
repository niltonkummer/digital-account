package user

import (
	"digital-account/application/models"
	ginJwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func (l *User) Authenticator(c *gin.Context) (interface{}, error) {

	req := LoginRequest{}
	err := c.BindJSON(&req)
	if err != nil {
		return nil, ginJwt.ErrMissingLoginValues
	}

	user, err := l.app.Repository.Login().Auth(c, req.CPF, "")
	if err != nil {
		return nil, ginJwt.ErrFailedAuthentication
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Secret), []byte(req.Secret)); err != nil {
		return nil, ginJwt.ErrFailedAuthentication
	}

	return user, nil
}

func (l *User) Authorizer(data interface{}, c *gin.Context) bool {

	if _, ok := data.(*models.User); !ok {
		return false
	}

	return true
}
