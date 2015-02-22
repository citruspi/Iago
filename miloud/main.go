package main

import (
	conf "github.com/citruspi/Milou/configuration"
	"github.com/citruspi/Milou/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	conf.Process()

	router := gin.Default()

	router.POST("/", handlers.Notification)

	router.Run(conf.Web.Address)
}
