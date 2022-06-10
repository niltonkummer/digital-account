package application

import (
	"digital-account/application/api/setup"
	"digital-account/application/config"
	"digital-account/application/models"
	"digital-account/application/repository"
	"flag"
	"log"
	"os"
	"strings"

	"gorm.io/driver/sqlite"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/rs/zerolog"
)

var (
	port       = flag.String("port", ":8080", "server port")
	configPath = flag.String("config", "settings/debug.yaml", "config file")
	listener   []string
	app        *config.App
)

type envType string

const (
	EnvDebug   envType = "debug"
	EnvRelease envType = "release"
)

func init() {

	flag.Parse()
	if *port != "" {
		listener = []string{*port}
	}

}

func setupDB(env config.Environment, dsn string) (*gorm.DB, error) {
	var dialect gorm.Dialector

	debug := false

	switch env {

	case config.ReleaseEnvironment:
		dialect = postgres.Open(dsn)
	default:
		fallthrough
	case config.DebugEnvironment:
		debug = true
		dialect = sqlite.Open(dsn)
	}
	db, err := gorm.Open(dialect)
	if err != nil {
		return nil, err
	}
	if debug {
		db = db.Debug()
	}
	return db, nil
}

func setupRepository(a *config.App) error {

	err := a.DB.AutoMigrate(&models.Account{}, &models.User{}, &models.Transfer{})
	if err != nil {
		return err
	}

	a.Repository = repository.Config(a.DB)

	return nil
}

func Run() {

	env := func() config.Environment {
		e := strings.ToLower(os.Getenv("ENVIRONMENT"))
		switch e {
		case string(EnvRelease):
			return config.ReleaseEnvironment
		default:
			return config.DebugEnvironment
		}
	}()

	settings, err := config.New(*configPath)
	if err != nil {
		log.Fatal("could not load settings: ", err)
	}
	app = &config.App{
		Logger:   zerolog.New(os.Stdout),
		Settings: settings,
	}

	app.DB = func() *gorm.DB {
		dbDialect := os.Getenv("DB")
		if dbDialect == "" {
			log.Fatal("$DB must be set")
		}
		db, err := setupDB(env, dbDialect)
		if err != nil {
			app.Logger.Fatal().Err(err).Msg("setupDB")
		}
		return db
	}()

	err = setupRepository(app)
	if err != nil {
		app.Logger.Fatal().Err(err).Msg("setupRepository")
	}

	setup.ConfigRoutes(app)

}
