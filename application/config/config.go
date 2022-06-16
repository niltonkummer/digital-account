package config

import (
	"digital-account/application/repository"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type Environment int

const (
	TestEnvironment    Environment = 0
	DebugEnvironment   Environment = 1
	ReleaseEnvironment Environment = 2
)

var envString = map[Environment]string{
	TestEnvironment:    "test",
	DebugEnvironment:   "debug",
	ReleaseEnvironment: "release",
}

func (e Environment) String() string {
	return envString[e]
}

type App struct {
	Router      gin.IRouter
	DB          *gorm.DB
	Logger      zerolog.Logger
	Environment Environment
	Settings    Settings
	Repository  repository.Repository
}

func New(path string) (Settings, error) {
	v := viper.New()

	v.SetConfigFile(path)
	v.AddConfigPath(".")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	err := v.ReadInConfig()
	if err != nil {
		return nil, err
	}

	s := &settings{
		source: v,
	}

	return s, nil
}

type Settings interface {
	Get(string) interface{}
	String(string) string
	Float64(string) float64
	Int(string) int
	Int64(string) int64
	Bool(string) bool
	Strings(string) []string
	Set(string, interface{})
}

type settings struct {
	source *viper.Viper
}

func (s *settings) get(key string) interface{} {
	value := s.source.Get(key)

	return value
}

func (s *settings) set(key string, value interface{}) {
	s.source.Set(key, value)
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

func (s *settings) Set(key string, value interface{}) {
	s.set(key, value)
}
