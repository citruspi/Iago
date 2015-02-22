package configuration

import (
	"log"

	"github.com/FogCreek/mini"
)

func Process() *mini.Config {
	config, err := mini.LoadConfiguration("iago.ini")

	if err != nil {
		log.Fatal(err)
	}

	return config
}
