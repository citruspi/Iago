package configuration

import (
	"flag"
	"log"

	"github.com/FogCreek/mini"
)

type webConfiguration struct {
	Address string
	Status  bool
}

type travisConfiguration struct {
	Authenticate bool
	Token        string
}

var (
	Web    webConfiguration
	Travis travisConfiguration
)

func Process() {
	path := flag.String("config", "/etc/iagod.ini", "Configuration file path")
	flag.Parse()

	config, err := mini.LoadConfiguration(*path)

	if err != nil {
		log.Fatal(err)
	}

	Web.Address = config.StringFromSection("Web", "Address", "127.0.01:8000")
	Web.Status = config.BooleanFromSection("Web", "Status", true)

	Travis.Authenticate = config.BooleanFromSection("Travis", "Authenticate", false)
	Travis.Token = config.StringFromSection("Travis", "Token", "")
}
