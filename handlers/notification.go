package handlers

import (
	"github.com/citruspi/iago/notification"
	"github.com/citruspi/milou/projects"
	"github.com/gin-gonic/gin"
)

func Notification(c *gin.Context) {
	var n notification.Notification

	c.Bind(&n)

	go projects.Process(n)
	c.JSON(200, "")
}
