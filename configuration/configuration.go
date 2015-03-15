package configuration

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
)

type Configuration struct {
	Mode string `json:"mode"`
	Web  struct {
		Address string `json:"address"`
	} `json:"web"`
	TravisCI struct {
		Authenticate bool   `json:"authenticate"`
		Token        string `json:"token"`
	} `json:"travis-ci"`
	Redis struct {
		Address string `json:"timeout"`
		Timeout int    `json:"timeout"`
	} `json:"redis"`
	Projects string
}

var (
	path *string
)

func init() {
	path = flag.String("config", "/etc/milou.conf", "Configuration file path")
	flag.Parse()
}

func Load() Configuration {
	var conf Configuration

	source, err := ioutil.ReadFile(*path)

	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(source, &conf)

	if err != nil {
		log.Fatal(err)
	}

	return conf
}
