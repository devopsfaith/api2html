package generator

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type Scanner interface {
	Scan() []TmplFolder
}

type TmplFolder struct {
	Path    string
	Content []string
}

func NewScanner(ts []string) Scanner {
	return TmplScanner(ts)
}

type TmplScanner []string

func (ts TmplScanner) Scan() []TmplFolder {
	res := []TmplFolder{}
	for _, prefix := range ts {
		if tmpls := getTemplatesInFolder(prefix); len(tmpls) > 0 {
			res = append(res, TmplFolder{
				Content: tmpls,
				Path:    prefix,
			})
		}
	}
	return res
}

func getTemplatesInFolder(prefix string) []string {
	templates := []string{}
	for _, fileName := range []string{"config.json", "Dockerfile"} {
		if _, err := os.Stat(fmt.Sprintf("%s/%s", prefix, fileName)); os.IsNotExist(err) {
			log.Println(err.Error())
			continue
		}
		templates = append(templates, fileName)
	}
	for _, folder := range []string{"tmpl", "static"} {
		files, err := ioutil.ReadDir(fmt.Sprintf("%s/%s", prefix, folder))
		if err != nil {
			log.Println(err.Error())
			continue
		}

		for _, f := range files {
			templates = append(templates, fmt.Sprintf("%s/%s", folder, f.Name()))
		}
	}
	return templates
}
