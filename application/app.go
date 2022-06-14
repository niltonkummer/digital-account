package application

import (
	apiConfig "digital-account/application/api/config"
	"digital-account/application/config"
	"digital-account/application/db"
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
	configPath = flag.String("config", "settings/debug.yaml", "config file")
	app        *config.App
)

type envType string

const (
	EnvTest    envType = "test"
	EnvDebug   envType = "debug"
	EnvRelease envType = "release"
)

func init() {
	flag.Parse()
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
	d, err := gorm.Open(dialect)
	if err != nil {
		return nil, err
	}
	if debug {
		d = d.Debug()
	}
	return d, nil
}

func setupRepository(a *config.App) error {

	err := db.Setup(a.DB)
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
		case string(EnvTest):
			return config.TestEnvironment
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

		dbDialect := app.Settings.String("db.dialect")

		d, err := setupDB(env, dbDialect)
		if err != nil {
			app.Logger.Fatal().Err(err).Msg("setupDB")
		}
		return d
	}()

	err = setupRepository(app)
	if err != nil {
		app.Logger.Fatal().Err(err).Msg("setupRepository")
	}

	apiConfig.Routes(app)

}
