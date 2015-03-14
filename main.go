package main

import (
	conf "github.com/citruspi/iago/configuration"
	"github.com/citruspi/iago/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	conf.Process()

	router := gin.Default()

	router.POST("/announce/", handlers.Announce)

	if conf.Web.Status {
		router.GET("/status/", handlers.Status)
	}

	router.Run(conf.Web.Address)
}
