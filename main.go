package main

import (
	"github.com/gin-gonic/gin"
)

type Host struct {
	Hostname string `json:"hostname"`
}

var (
	hosts []Host
)

func main() {
	router := gin.Default()

	router.POST("/checkin/", func(c *gin.Context) {
		var host Host

		c.Bind(&host)

		if host.Hostname != "" {
			exists := false
			for _, h := range hosts {
				if h.Hostname == host.Hostname {
					exists = true
				}
			}
			if !exists {
				hosts = append(hosts, host)
			}
		}

		c.String(200, "")
	})

	router.Run(":8080")
}
