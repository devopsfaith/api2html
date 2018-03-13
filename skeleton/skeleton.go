package skeleton

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	// Import the statikFS from the generated files
	_ "github.com/devopsfaith/api2html/statik"
	"github.com/rakyll/statik/fs"
)

var (
	blogContents = []string{
		"/blog/config/es_ES/config.ini",
		"/blog/config/es_ES/routes.ini",
		"/blog/config/global/config.ini",
		"/blog/config/global/routes.ini",
		"/blog/i18n/es_ES.ini",
		"/blog/i18n/en_US.ini",
		"/blog/sources/es_ES/static/404",
		"/blog/sources/es_ES/static/500",
		"/blog/sources/es_ES/tmpl/home.mustache",
		"/blog/sources/global/config.json",
		"/blog/sources/global/Dockerfile",
		"/blog/sources/global/static/404",
		"/blog/sources/global/static/500",
		"/blog/sources/global/static/hello.txt",
		"/blog/sources/global/static/robots.txt",
		"/blog/sources/global/static/sitemap.xml",
		"/blog/sources/global/tmpl/home.mustache",
		"/blog/sources/global/tmpl/main_layout.mustache",
		"/blog/sources/global/tmpl/post.mustache",
	}
)

// New returns a statikFS Skel
func New(outputPath string, fileList []string) Skel {
	return &statikSkel{outputPath: outputPath, fileList: fileList}
}

// NewBlog returns a statikSkel with the blog example contents
func NewBlog(outputPath string) Skel {
	return &statikSkel{outputPath: outputPath, fileList: blogContents}
}

// Skel defines the interface for creating skeleton files
type Skel interface {
	Create() error
}

type statikSkel struct {
	outputPath string
	fileList   []string
}

// Generate the skel from the statikFS
func (s *statikSkel) Create() error {
	statikFS, err := fs.New()
	if err != nil {
		return err
	}

	for _, name := range s.fileList {
		f, err := statikFS.Open(name)
		if err != nil {
			fmt.Printf("opening file %s: %s\n", name, err.Error())
			return err
		}
		buff := new(bytes.Buffer)
		_, err = buff.ReadFrom(f)
		f.Close()
		if err != nil {
			return err
		}
		path := filepath.Join(s.outputPath, filepath.Dir(name))
		if _, err := os.Stat(path); os.IsNotExist(err) {
			err = os.MkdirAll(path, os.ModePerm)
			if err != nil {
				return err
			}
		}
		filename := filepath.Join(s.outputPath, name)
		err = ioutil.WriteFile(filename, buff.Bytes(), os.ModePerm)
		if err != nil {
			return err
		}
		fmt.Printf("Creating skeleton file: %s\n", filename)
	}
	return nil
}
