package engine

import (
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/cbroglie/mustache"
)

func NewMustacheRendererMap(cfg Config) (map[string]*MustacheRenderer, error) {
	result := map[string]*MustacheRenderer{}
	for _, section := range []map[string]string{cfg.Templates, cfg.Layouts} {
		for name, path := range section {
			templateFile, err := os.Open(path)
			if err != nil {
				log.Println("reading", path, ":", err.Error())
				return result, err
			}
			renderer, err := NewMustacheRenderer(templateFile)
			templateFile.Close()
			if err != nil {
				log.Println("parsing", path, ":", err.Error())
				return result, err
			}
			result[name] = renderer
		}
	}
	return result, nil
}

func NewMustacheRenderer(r io.Reader) (*MustacheRenderer, error) {
	tmpl, err := newMustacheTemplate(r)
	if err != nil {
		return nil, err
	}
	return &MustacheRenderer{tmpl}, nil
}

type MustacheRenderer struct {
	tmpl *mustache.Template
}

func (m MustacheRenderer) Render(w io.Writer, v interface{}) error {
	return m.tmpl.FRender(w, v)
}

func NewLayoutMustacheRenderer(t, l io.Reader) (*LayoutMustacheRenderer, error) {
	tmpl, err := newMustacheTemplate(t)
	if err != nil {
		return nil, err
	}
	layout, err := newMustacheTemplate(l)
	if err != nil {
		return nil, err
	}
	return &LayoutMustacheRenderer{tmpl, layout}, nil
}

type LayoutMustacheRenderer struct {
	tmpl   *mustache.Template
	layout *mustache.Template
}

func (m LayoutMustacheRenderer) Render(w io.Writer, v interface{}) error {
	return m.tmpl.FRenderInLayout(w, m.layout, v)
}

func newMustacheTemplate(r io.Reader) (*mustache.Template, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return mustache.ParseString(string(data))
}
