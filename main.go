package main

import (
	"net/http"

	"github.com/citruspi/milou/configuration"
	"github.com/citruspi/milou/handlers"
	"github.com/citruspi/milou/projects"
)

var (
	conf configuration.Configuration
)

func init() {
	conf = configuration.Load()
}

func main() {
	if conf.Mode == "server" || conf.Mode == "standalone" {
		http.HandleFunc("/webhooks/travis/", handlers.Travis)
		http.ListenAndServe(conf.Web.Address, nil)
	} else if conf.Mode == "client" {
		projects.Subscribe()
	}
}
