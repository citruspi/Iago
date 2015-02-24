package handlers

import (
	"time"

	"github.com/citruspi/iago/host"
	"github.com/gin-gonic/gin"
)

type status struct {
	Hosts  []host.Host `json:"hosts"`
	Uptime float64     `json:"uptime"`
}

var (
	boot = time.Now()
)

func Status(c *gin.Context) {
	s := status{}
	s.Hosts = host.List
	s.Uptime = time.Now().Sub(boot).Seconds()

	c.JSON(200, s)
}
