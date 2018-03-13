package skeleton

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	// Import the statikFS from the generated files
	_ "github.com/devopsfaith/api2html/statik"
	"github.com/rakyll/statik/fs"
)

// New returns a statikFS Skel
func New(outputPath string) Skel {
	return &statikSkel{outputPath: outputPath}
}

// Skel defines the interface for creating skeleton files
type Skel interface {
	Create() error
}

type statikSkel struct {
	outputPath string
}

// Generate the skel from the statikFS
func (s *statikSkel) Create() error {
	statikFS, err := fs.New()
	if err != nil {
		return err
	}

	if _, err := os.Stat(s.outputPath); err == nil {
		return errors.New("output directory already exists")
	}

	for _, name := range []string{
		"/config/es_ES/config.ini",
		"/config/es_ES/routes.ini",
		"/config/global/config.ini",
		"/config/global/routes.ini",
		"/i18n/es_ES.ini",
		"/i18n/en_US.ini",
		"/sources/es_ES/static/404",
		"/sources/es_ES/static/500",
		"/sources/es_ES/tmpl/home.mustache",
		"/sources/global/config.json",
		"/sources/global/Dockerfile",
		"/sources/global/static/404",
		"/sources/global/static/500",
		"/sources/global/static/hello.txt",
		"/sources/global/static/robots.txt",
		"/sources/global/static/sitemap.xml",
		"/sources/global/tmpl/home.mustache",
		"/sources/global/tmpl/main_layout.mustache",
		"/sources/global/tmpl/post.mustache",
	} {
		f, err := statikFS.Open(name)
		if err != nil {
			fmt.Printf("opening file %s: %s\n", name, err.Error())
			return err
		}
		defer f.Close()
		buff := new(bytes.Buffer)
		_, err = buff.ReadFrom(f)
		if err != nil {
			return err
		}
		path := fmt.Sprintf("%s/%s", s.outputPath, filepath.Dir(name))
		if _, err := os.Stat(path); os.IsNotExist(err) {
			err = os.MkdirAll(path, os.ModePerm)
			if err != nil {
				return err
			}
		}

		err = ioutil.WriteFile(fmt.Sprintf("%s%s", s.outputPath, name), buff.Bytes(), os.ModePerm)
		if err != nil {
			return err
		}
		fmt.Printf("Creating skeleton file: %s%s\n", s.outputPath, name)
	}
	return nil
}
