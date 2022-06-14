package config

import (
	"digital-account/application/api/accounts"
	"digital-account/application/api/common"
	"digital-account/application/api/transfers"
	"digital-account/application/api/users"
	"digital-account/application/config"
	"digital-account/application/models"
	"strings"

	"time"

	"log"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

func Routes(app *config.App) {

	gin.SetMode(app.Environment.String())

	router := gin.Default()

	app.Router = router

	loginH := users.Config(app)
	authMiddleware := Auth(app, loginH)

	router.POST("/login", authMiddleware.LoginHandler)

	api := router.Group("/api")
	{
		apiAccounts := accounts.Config(app)
		notAuthApi := api.Group("")
		{
			accountsHandler := notAuthApi.Group("/accounts")
			accountsHandler.POST("", apiAccounts.CreateHandler)
		}

		authApi := api.Group("")
		{
			authApi.Use(authMiddleware.MiddlewareFunc())

			accountsHandler := authApi.Group("/accounts")
			{
				accountsHandler.GET("", apiAccounts.ListHandler)
				accountsHandler.GET("/:account_id/balance", apiAccounts.BalanceHandler)

			}

			transfersHandler := authApi.Group("/transfers")
			{
				apiTransfers := transfers.Config(app)
				transfersHandler.POST("", apiTransfers.CreateHandler)
				transfersHandler.GET("", apiTransfers.ListHandler)
			}
		}
	}

	app.Logger.Info().Msgf("Starting HTTP on %v", app.Settings.String("container.port"))
	err := router.Run(app.Settings.String("container.port"))
	if err != nil {
		app.Logger.Fatal().Err(err).Msg("[HTTP LISTEN]")
	}
}

func Auth(app *config.App, l *users.User) *jwt.GinJWTMiddleware {

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
			user, _ := app.Repository.User().Get(c, int64(claims[common.IdentityKey].(float64)))
			return user
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
