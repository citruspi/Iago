package main

import (
	conf "github.com/citruspi/Milou/configuration"
	"github.com/citruspi/Milou/handlers"
	"github.com/citruspi/Milou/iago"
	"github.com/citruspi/Milou/projects"
	"github.com/gin-gonic/gin"
)

func main() {
	conf.Process()

	for _, project := range projects.List {
		project.Deploy()
	}

	go iago.CheckIn()

	router := gin.Default()

	router.POST("/", handlers.Notification)

	router.Run(conf.Web.Address)
}
