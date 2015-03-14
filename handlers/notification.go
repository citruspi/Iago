package handlers

import (
	"github.com/citruspi/iago/notifications"
	"github.com/citruspi/milou/projects"
	"github.com/gin-gonic/gin"
)

func Notification(c *gin.Context) {
	var n notifications.Notification

	c.Bind(&n)

	go projects.Process(n)
	c.JSON(200, "")
}
