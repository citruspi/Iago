package main

import (
	conf "github.com/citruspi/milou/configuration"
	"github.com/citruspi/milou/handlers"
	"github.com/citruspi/milou/iago"
	"github.com/citruspi/milou/projects"
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
