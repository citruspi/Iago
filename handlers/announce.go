package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

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

	payload, _ := json.Marshal(notification.Payload)
	body := bytes.NewBuffer(payload)

	for _, host := range host.List {
		urlStr := host.URL()

		req, err := http.NewRequest("POST", urlStr, body)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)

		if err != nil {
			log.Fatal(err)
		}

		defer resp.Body.Close()
	}

	c.JSON(200, "")
}
