package config

import (
	"digital-account/application/repository"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type Environment int

const (
	DebugEnvironment   Environment = 1
	ReleaseEnvironment Environment = 2
)

var envString = map[Environment]string{
	DebugEnvironment:   "debug",
	ReleaseEnvironment: "release",
}

func (e Environment) String() string {
	return envString[e]
}

type Config struct {
	HttpListen string
	DBDialect  gorm.Dialector
}

type Settings interface {
	Get(string) interface{}
	String(string) string
	Float64(string) float64
	Int(string) int
	Int64(string) int64
	Bool(string) bool
	Strings(string) []string
}

type App struct {

	// Config
	conf *Config

	// App
	Router gin.IRouter
	DB     *gorm.DB
	Logger zerolog.Logger

	Environment Environment

	Settings Settings

	// Repository
	Repository repository.Repository
}

func New(path string) (Settings, error) {
	v := viper.New()

	v.SetConfigFile(path)

	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	s := &settings{
		source: v,
	}

	return s, nil
}

type settings struct {
	source *viper.Viper
}

func (s *settings) get(key string) interface{} {
	value := s.source.Get(key)

	return value
}

func (s *settings) Get(key string) interface{} {
	value := s.source.Get(key)

	return value
}

func (s *settings) String(key string) string {
	return cast.ToString(s.get(key))
}

func (s *settings) Bool(key string) bool {
	return cast.ToBool(s.get(key))
}

func (s *settings) Int(key string) int {
	return cast.ToInt(s.get(key))
}

func (s *settings) Int64(key string) int64 {
	return cast.ToInt64(s.get(key))
}

func (s *settings) Float64(key string) float64 {
	return cast.ToFloat64(s.get(key))
}

func (s *settings) Strings(key string) []string {
	return cast.ToStringSlice(s.get(key))
}
