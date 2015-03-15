package main

import (
	"net/http"

	conf "github.com/citruspi/iago/configuration"
	"github.com/citruspi/iago/handlers"
)

func main() {
	conf.Process()

	http.HandleFunc("/", handlers.TravisWebhook)
	http.ListenAndServe(conf.Web.Address, nil)
}
