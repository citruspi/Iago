package handlers

import (
	log "github.com/Sirupsen/logrus"
	"github.com/citruspi/iago/host"
	"github.com/gin-gonic/gin"
)

func init() {
	log.SetLevel(log.DebugLevel)
}

func CheckIn(c *gin.Context) {
	var h host.Host

	c.Bind(&h)

	log.Info("Received host checkin")

	if h.Hostname != "" {
		h.CheckIn()

		c.JSON(200, "")
	} else {
		log.Error("Host is incomplete")

		var response struct {
			Message string `json:"message"`
		}

		response.Message = "Invalid hostname."
		c.JSON(400, response)
	}
}
