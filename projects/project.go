package projects

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	log "github.com/Sirupsen/logrus"
	"github.com/citruspi/milou/configuration"
	"github.com/citruspi/milou/notifications"
)

type Project struct {
	Owner      string `json:"owner"`
	Repository string `json:"repository"`
	Version    struct {
		Type  string `json:"type"`
		Value string `json:"value"`
		Regex bool   `json:"regex"`
	} `json:"version"`
	Identifier string `json:"identifier"`
	Path       string `json:"path"`
}

var (
	conf configuration.Configuration
	list []Project
)

func init() {
	conf = configuration.Load()

	files, err := ioutil.ReadDir(conf.Projects)

	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		var project Project

		source, err := ioutil.ReadFile(conf.Projects + file.Name())

		if err != nil {
			log.Fatal(err)
		}

		err = json.Unmarshal(source, &project)

		if err != nil {
			log.Fatal(err)
		}

		list = append(list, project)
	}

	if conf.Mode == "standalone" || conf.Mode == "client" {
		DeployAll()
	}
}

func (p Project) BasePath() string {
	buffer := new(bytes.Buffer)

	t, _ := template.New("basepath").Parse(p.Path)
	t.Execute(buffer, p)

	return string(buffer.Bytes())
}

func (p Project) ArchivePath() string {
	var buffer bytes.Buffer

	buffer.WriteString(p.TemporaryPath())
	buffer.WriteString(p.Version.Value)
	buffer.WriteString(".zip")

	return string(buffer.Bytes())
}

func (p Project) TemporaryPath() string {
	var buffer bytes.Buffer

	buffer.WriteString(p.BasePath()[:len(p.BasePath())-1])
	buffer.WriteString(".milou/")

	return string(buffer.Bytes())
}

func (p Project) ExtractPath() string {
	var buffer bytes.Buffer

	buffer.WriteString(p.TemporaryPath())
	buffer.WriteString(p.Repository)

	return string(buffer.Bytes())
}

func (p Project) ArchiveLocation() string {
	var buffer bytes.Buffer

	buffer.WriteString("https://s3.amazonaws.com/")
	buffer.WriteString(p.Identifier)
	buffer.WriteString("/")
	buffer.WriteString(p.Version.Value)
	buffer.WriteString(".zip")

	return string(buffer.Bytes())
}

func (p Project) Extract() error {
	err := os.MkdirAll(p.ExtractPath(), 0700)

	if err != nil {
		return err
	}

	r, err := zip.OpenReader(p.ArchivePath())
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		fpath := filepath.Join(p.ExtractPath(), f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, f.Mode())
		} else {
			var fdir string
			if lastIndex := strings.LastIndex(fpath, string(os.PathSeparator)); lastIndex > -1 {
				fdir = fpath[:lastIndex]
			}

			err = os.MkdirAll(fdir, f.Mode())
			if err != nil {
				log.Fatal(err)
				return err
			}
			f, err := os.OpenFile(
				fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer f.Close()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (p Project) Download() {
	response, err := http.Get(p.ArchiveLocation())
	defer response.Body.Close()

	if err != nil {
		log.Fatal(err)
	}

	archive, err := os.Create(p.ArchivePath())
	defer archive.Close()

	_, err = io.Copy(archive, response.Body)

	if err != nil {
		log.Fatal(err)
	}
}

func (p Project) Place() {
	err := os.RemoveAll(p.BasePath())

	if err != nil {
		log.Fatal(err)
	}

	err = os.Rename(p.ExtractPath(), p.BasePath())

	if err != nil {
		log.Fatal(err)
	}
}

func (p Project) Prepare() {
	if _, err := os.Stat(p.TemporaryPath()); os.IsNotExist(err) {
		err = os.MkdirAll(p.TemporaryPath(), 0700)

		if err != nil {
			log.Fatal(err)
		}
	}
}

func (p Project) CleanUp() {
	err := os.RemoveAll(p.TemporaryPath())

	if err != nil {
		log.Fatal(err)
	}
}

func (p Project) Deploy() {
	p.Prepare()
	p.Download()
	p.Extract()
	p.Place()
	p.CleanUp()
}

func DeployAll() {
	for _, project := range list {
		project.Deploy()
	}
}

func Process(n notifications.Notification) {
	for _, project := range list {
		if project.Repository != n.Repository {
			continue
		}

		if project.Owner != n.Owner {
			continue
		}

		if project.Version.Type == "commit" {
			if project.Version.Value == n.Commit {
				project.Deploy()
				continue
			}
		}

		if project.Version.Type == "branch" {
			if project.Version.Regex {
				match, _ := regexp.MatchString(project.Version.Value, n.Branch)

				if match {
					project.Deploy()
				}
			} else if project.Version.Value == n.Branch {
				project.Deploy()
			}
		}
	}
}
