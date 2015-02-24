package handlers

import (
	"github.com/citruspi/iago/notification"
	conf "github.com/citruspi/milou/configuration"
	"github.com/citruspi/milou/projects"
	"github.com/gin-gonic/gin"
)

func Notification(c *gin.Context) {
	var n notification.Notification

	c.Bind(&n)

	if conf.Notification.Signed {
		if n.Verify(conf.Notification.PublicKey) {
			go projects.Process(n)
			c.JSON(200, "")
		} else {
			c.JSON(400, "")
		}
	} else {
		c.JSON(200, "")
	}
}
