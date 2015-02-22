package handlers

import (
	"github.com/citruspi/Iago/host"
	"github.com/gin-gonic/gin"
)

func CheckIn(c *gin.Context) {
	var h host.Host

	c.Bind(&h)

	if h.Hostname != "" {
		h.CheckIn()

		c.JSON(200, "")
	} else {
		var response struct {
			Message string `json:"message"`
		}

		response.Message = "Invalid hostname."
		c.JSON(400, response)
	}
}
