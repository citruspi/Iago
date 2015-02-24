package projects

import (
	"archive/zip"
	"bytes"
	"github.com/citruspi/iago/notification"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Project struct {
	Name       string
	Owner      string
	Repository string
	Version    string
	Identifier string
	Domain     string
	Subdomain  string
	Type       string
}

var (
	List []Project
)

func (p Project) Path() string {
	var buffer bytes.Buffer

	if p.Type == "static" {
		buffer.WriteString("/srv/")
		buffer.WriteString(p.Domain)
		buffer.WriteString("/")
		buffer.WriteString(p.Subdomain)
		buffer.WriteString("/")
	}

	return string(buffer.Bytes())
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

	buffer.WriteString(p.Path()[:len(p.Path())-1])
	buffer.WriteString(".milou/")

	return string(buffer.Bytes())
}

func (p Project) ExtractPath() string {
	var buffer bytes.Buffer

	buffer.WriteString(p.TemporaryPath())
	buffer.WriteString(p.Subdomain)

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
	err := os.RemoveAll(p.Path())

	if err != nil {
		log.Fatal(err)
	}

	err = os.Rename(p.ExtractPath(), p.Path())

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

func Process(n notification.Notification) {
	for _, project := range List {
		if project.Repository == n.Repository {
			if project.Owner == n.Owner {
				if string(project.Version[0]) == "#" {
					if project.Version[:len(project.Version)-1] == n.Commit {
						project.Deploy()
					}
				} else if string(project.Version[0]) == "@" {
					if project.Version[:len(project.Version)-1] == n.Branch {
						project.Deploy()
					}
				}
			}
		}
	}
}
