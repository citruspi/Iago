package handlers

import (
	"github.com/citruspi/Iago/host"
	"github.com/citruspi/Iago/travis"
	"github.com/gin-gonic/gin"
)

func Announce(c *gin.Context) {
	var announcement travis.Announcement

	c.Bind(&announcement)

	if authorization, ok := c.Request.Header["Authorization"]; ok {
		announcement.Authorization = authorization[0]
	}

	if !announcement.Valid() {
		c.JSON(400, "")
		return
	}

	go host.Notify(announcement)

	c.JSON(200, "")
}
