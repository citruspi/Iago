package main

import (
	"time"

	"github.com/gin-gonic/gin"
)

type Host struct {
	Hostname string    `json:"hostname"`
	TTL      time.Time `json:"ttl"`
}

var (
	hosts []Host
	TTL   = 30
)

func buildTTL() time.Duration {
	expiration = time.Now().UTC()
	expiration.Add(time.Duration(TTL) * time.Second)

	return expiration
}

func main() {
	router := gin.Default()

	router.POST("/checkin/", func(c *gin.Context) {
		var host Host

		c.Bind(&host)

		if host.Hostname != "" {
			exists := false

			for i, h := range hosts {
				if h.Hostname == host.Hostname {
					exists = true
					hosts[i].TTL = buildTTL()
				}
			}

			if !exists {
				host.TTL = buildTTL()
				hosts = append(hosts, host)
			}

			c.JSON(200, "")
		} else {
			var response struct {
				Message string `json:"message"`
			}
			response.Message = "Invalid hostname."
			c.JSON(400, response)
		}
	})

	router.GET("/status/", func(c *gin.Context) {
		c.JSON(200, hosts)
	})

	router.Run(":8080")
}
