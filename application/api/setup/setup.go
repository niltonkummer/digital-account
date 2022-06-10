package setup

import (
	"digital-account/application/api/accounts"
	"digital-account/application/api/common"
	"digital-account/application/api/user"
	"digital-account/application/config"
	"digital-account/application/models"
	"strings"

	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"log"
)

func ConfigRoutes(app *config.App) {

	gin.SetMode(app.Environment.String())

	router := gin.Default()

	app.Router = router

	loginH := user.Config(app)
	authMiddleware := ConfigAuth(app, loginH)

	loginR := router.Group("/login")
	{

		loginR.POST("", authMiddleware.LoginHandler)
	}

	api := router.Group("/api")
	{

		api.Use(authMiddleware.MiddlewareFunc())

		accountsHandler := api.Group("/accounts")
		{
			accHandler := accounts.Config(app)
			accountsHandler.GET("", accHandler.ListHandler)
			accountsHandler.GET("/:account_id/balance", accHandler.GetHandler)
			accountsHandler.POST("", accHandler.CreateHandler)
		}
	}

	err := router.Run(app.Settings.String("container.port"))
	if err != nil {
		log.Fatalln(err)
	}
}

func ConfigAuth(app *config.App, l *user.User) *jwt.GinJWTMiddleware {

	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm: app.Settings.String("auth.realm"),
		Key:   []byte(app.Settings.String("auth.key")),
		Timeout: func() time.Duration {
			if d, err := time.ParseDuration(app.Settings.String("auth.timeout")); err != nil {
				return time.Hour
			} else {
				return d
			}
		}(),
		MaxRefresh:  time.Hour,
		IdentityKey: common.IdentityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*models.User); ok {
				return jwt.MapClaims{
					common.IdentityKey: v.ID,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &models.User{
				Model: models.Model{
					ID: int64(claims[common.IdentityKey].(float64)),
				},
			}
		},
		Authenticator: l.Authenticator,
		Authorizator:  l.Authorizer,
		TokenLookup: strings.Join(
			[]string{
				"header: Authorization",
				"query: token",
				"cookie: jwt",
			}, ","),
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	errInit := authMiddleware.MiddlewareInit()
	if errInit != nil {
		log.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
	}

	return authMiddleware
}
