package main

import (
	"digital-account/application"

	_ "github.com/heroku/x/hmetrics/onload"
)

func main() {
	application.Run()
}
