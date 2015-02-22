package main

import (
	conf "github.com/citruspi/Iago/configuration"
	"github.com/citruspi/Iago/handlers"
	"github.com/citruspi/Iago/host"
	"github.com/gin-gonic/gin"
)

func main() {
	conf.Process()

	go host.Cleanup()

	router := gin.Default()

	router.POST("/checkin/", handlers.CheckIn)
	router.POST("/announce/", handlers.Announce)

	if conf.Web.Status {
		router.GET("/status/", handlers.Status)
	}

	router.Run(conf.Web.Address)
}
