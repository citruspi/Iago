package projects

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/citruspi/milou/configuration"
	"github.com/citruspi/milou/notifications"
)

type Project struct {
	Owner      string
	Repository string
	Version    string
	Identifier string
	Path       string
}

var (
	conf     configuration.Configuration
	projects []Project
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

		projects = append(projects, project)
	}
}

func (p Project) ArchivePath() string {
	var buffer bytes.Buffer

	buffer.WriteString(p.TemporaryPath())
	buffer.WriteString(p.Version)
	buffer.WriteString(".zip")

	return string(buffer.Bytes())
}

func (p Project) TemporaryPath() string {
	var buffer bytes.Buffer

	buffer.WriteString(p.Path[:len(p.Path)-1])
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
	buffer.WriteString(p.Version)
	buffer.WriteString(".zip")

	return string(buffer.Bytes())
}

func (p Project) Extract() error {
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
	err := os.RemoveAll(p.Path)

	if err != nil {
		log.Fatal(err)
	}

	err = os.Rename(p.ExtractPath(), p.Path)

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
	for _, project := range projects {
		project.Deploy()
	}
}

func Process(n notifications.Notification) {
	for _, project := range projects {
		if project.Repository == n.Repository {
			if project.Owner == n.Owner {
				if project.Version == n.Commit {
					project.Deploy()
				} else if project.Version == n.Branch {
					project.Deploy()
				}
			}
		}
	}
}
