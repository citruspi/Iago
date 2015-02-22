package main

import (
	"github.com/citruspi/Iago/configuration"
	"github.com/citruspi/Iago/handlers"
	"github.com/citruspi/Iago/host"
	"github.com/gin-gonic/gin"
)

var (
	conf configuration.Configuration
)

func main() {
	configuration.Process()
	conf = configuration.Conf

	go host.Cleanup()

	router := gin.Default()

	router.POST("/checkin/", handlers.CheckIn)
	router.POST("/announce/", handlers.Announce)

	if conf.Web.Status {
		router.GET("/status/", handlers.Status)
	}

	router.Run(conf.Web.Address)
}
