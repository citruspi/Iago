package handlers

import (
	"github.com/gin-gonic/gin"
)

func Notification(c *gin.Context) {
	c.String(200, "yolo")
}
