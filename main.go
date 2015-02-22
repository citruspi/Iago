package main

import (
	"time"

	"github.com/citruspi/Iago/configuration"
	"github.com/citruspi/Iago/handlers"
	"github.com/citruspi/Iago/host"
	"github.com/gin-gonic/gin"
)

var (
	conf configuration.Configuration
)

func cleanup() {
	for {
		for i, h := range host.List {
			if time.Now().UTC().After(h.Expiration) {
				diff := host.List
				diff = append(diff[:i], diff[i+1:]...)
				host.List = diff
			}
		}
		time.Sleep(time.Duration(conf.Host.TTL) * time.Second)
	}
}

func main() {
	configuration.Process()
	conf = configuration.Conf

	go cleanup()

	router := gin.Default()

	router.POST("/checkin/", handlers.CheckIn)
	router.POST("/announce/", handlers.Announce)

	if conf.Web.Status {
		router.GET("/status/", handlers.Status)
	}

	router.Run(conf.Web.Address)
}
