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

func buildTTL() time.Time {
	expiration := time.Now().UTC()
	expiration = expiration.Add(time.Duration(TTL) * time.Second)

	return expiration
}

func cleanup() {
	for {
		for i, h := range hosts {
			if time.Now().UTC().After(h.TTL) {
				diff := hosts
				diff = append(diff[:i], diff[i+1:]...)
				hosts = diff
			}
		}
		time.Sleep(time.Duration(TTL) * time.Second)
	}
}

func main() {
	go cleanup()

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
