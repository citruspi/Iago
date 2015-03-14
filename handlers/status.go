package handlers

import (
	"time"

	"github.com/gin-gonic/gin"
)

type status struct {
	Uptime float64 `json:"uptime"`
}

var (
	boot = time.Now()
)

func Status(c *gin.Context) {
	s := status{}
	s.Uptime = time.Now().Sub(boot).Seconds()

	c.JSON(200, s)
}
