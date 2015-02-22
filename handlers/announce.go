package handlers

import (
	"github.com/citruspi/Iago/host"
	"github.com/citruspi/Iago/travis"
	"github.com/gin-gonic/gin"
)

func Announce(c *gin.Context) {
	var notification travis.Notification

	c.Bind(&notification)

	if authorization, ok := c.Request.Header["Authorization"]; ok {
		notification.Authorization = authorization[0]
	}

	if !notification.Valid() {
		c.JSON(400, "")
		return
	}

	host.Notify(notification)

	c.JSON(200, "")
}
