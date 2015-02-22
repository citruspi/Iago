package configuration

import (
	"flag"
	"log"

	"github.com/FogCreek/mini"
)

func Process() {
	path := flag.String("config", "/etc/miloud.ini", "Configuration file path")
	flag.Parse()

	config, err := mini.LoadConfiguration(*path)

	if err != nil {
		log.Fatal(err)
	}
}
