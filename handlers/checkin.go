package handlers

import (
	"github.com/citruspi/Iago/host"
	"github.com/gin-gonic/gin"
)

func CheckIn(c *gin.Context) {
	var newHost host.Host

	c.Bind(&newHost)

	if newHost.Hostname != "" {
		for i, h := range host.List {
			if h.Hostname == newHost.Hostname {
				diff := host.List
				diff = append(diff[:i], diff[i+1:]...)
				host.List = diff
			}
		}

		newHost = newHost.Process()

		host.List = append(host.List, newHost)

		c.JSON(200, "")
	} else {
		var response struct {
			Message string `json:"message"`
		}
		response.Message = "Invalid hostname."
		c.JSON(400, response)
	}
}
