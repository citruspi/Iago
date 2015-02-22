package main

import (
	"github.com/citruspi/Iago/configuration"
	"github.com/citruspi/Iago/handlers"
	"github.com/citruspi/Iago/host"
	"github.com/gin-gonic/gin"
)

func main() {
	configuration.Process()

	go host.Cleanup()

	router := gin.Default()

	router.POST("/checkin/", handlers.CheckIn)
	router.POST("/announce/", handlers.Announce)

	if configuration.Web.Status {
		router.GET("/status/", handlers.Status)
	}

	router.Run(configuration.Web.Address)
}
