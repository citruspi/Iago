package configuration

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"

	"github.com/citruspi/milou/projects"
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
	Projects []struct {
		Owner      string `json:"owner"`
		Repository string `json:"repository"`
		Version    string `json:"version"`
		Identifier string `json:"identifier"`
		Path       string `json:"path"`
	} `json:"projects"`
}

var (
	path *string
)

func init() {
	path = flag.String("config", "/etc/milou.conf", "Configuration file path")
	flag.Parse()

	conf := Load()

	for _, project := range conf.Projects {
		projects.List = append(projects.List, projects.Project{
			Owner:      project.Owner,
			Repository: project.Repository,
			Version:    project.Version,
			Identifier: project.Identifier,
			Path:       project.Path,
		})
	}
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
