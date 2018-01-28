package generator

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/cbroglie/mustache"
)

type Renderer interface {
	Render(string, Data, Scanner) error
}

func NewRenderer(outputFolder string, r *regexp.Regexp) Renderer {
	return &Render{
		OutputFolder: outputFolder,
		Regexp:       r,
		Dumper:       mustacheRender,
	}
}

type Dumper func(string, string, Data) error

type Render struct {
	OutputFolder string
	Regexp       *regexp.Regexp
	Dumper       Dumper
}

func (r *Render) Render(iso string, data Data, scanner Scanner) error {
	if err := r.prepareOutputFolder(iso); err != nil {
		return err
	}
	for _, tmplFolder := range scanner.Scan() {
		if len(tmplFolder.Content) == 0 {
			log.Println("skipping the empty folder:", tmplFolder.Path)
			continue
		}
		for _, tmplName := range tmplFolder.Content {
			source := fmt.Sprintf("%s/%s", tmplFolder.Path, tmplName)
			target := fmt.Sprintf("%s/%s/%s", r.OutputFolder, iso, tmplName)

			if r.Regexp != nil && r.Regexp.Match([]byte(tmplName)) {
				log.Println("ignoring the source file:", source)
				continue
			}

			if err := r.Dumper(source, target, data); err != nil {
				log.Printf("rendering [%s] into [%s]: %s", source, target, err.Error())
				return err
			}
		}
		log.Println(tmplFolder.Path, "preprocessed")
	}
	return nil
}

func (r *Render) prepareOutputFolder(iso string) error {
	target := fmt.Sprintf("%s/%s", r.OutputFolder, iso)
	if err := os.RemoveAll(target); err != nil {
		return err
	}

	if err := os.MkdirAll(fmt.Sprintf("%s/tmpl", target), os.ModePerm); err != nil {
		return err
	}

	return os.MkdirAll(fmt.Sprintf("%s/static", target), os.ModePerm)
}

func mustacheRender(source, target string, data Data) error {
	if _, err := os.Stat(source); os.IsNotExist(err) {
		log.Println(source, "doesn't exist. Skipping.")
		return nil
	}

	log.Println("parsing the file", source)

	tmpl, err := mustache.ParseFile(source)
	if err != nil {
		return err
	}

	log.Println("writting the contents of", source, "into", target)

	file, err := os.OpenFile(target, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	buf := bytes.NewBuffer([]byte{})

	if err := tmpl.FRender(buf, data); err != nil {
		return err
	}

	file.Write(bytes.Replace(buf.Bytes(), []byte("&#34;"), []byte(`"`), -1))

	return nil
}
