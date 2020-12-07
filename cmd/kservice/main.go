package main

import (
	"net/http"

	"github.com/jinghzhu/kservice/pkg/api/v1/router"
	"github.com/jinghzhu/kservice/pkg/config"
	"github.com/jinghzhu/kservice/pkg/logger"
)

func init() {
	config.GetConfig()
	logger.SetLogLevel(logger.INFO)
}

func main() {
	g := config.GetConfig()
	logger.InfoFields("Start kservice", logger.Fields{"Config": g})
	r := router.DefaultRouter()
	logger.Error(http.ListenAndServe(g.ListenAddress, r))
}
