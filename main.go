package main

import (
	conf "github.com/citruspi/iago/configuration"
	"github.com/citruspi/iago/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	conf.Process()

	router := gin.Default()

	router.POST("/webhooks/travis/", handlers.TravisWebhook)

	router.Run(conf.Web.Address)
}
