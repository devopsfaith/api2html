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
	return mustache.ParseStringPartials(string(data), customPartialProvider)
}

type partialProvider struct {
	statics mustache.PartialProvider
	dynamc  mustache.PartialProvider
}

func (sp *partialProvider) Get(name string) (string, error) {
	if data, err := sp.statics.Get(name); err == nil {
		return data, nil
	}

	return sp.dynamc.Get(name)
}

var (
	customPartialProvider = &partialProvider{
		dynamc:  &mustache.FileProvider{},
		statics: &mustache.StaticProvider{Partials: map[string]string{"api2html/debug": debuggerTmpl}},
	}

	debuggerTmpl = `
<div>
	<h1>API2HTML Debugger</h1>
    <small>page generated at {{ Helper.Now }}</small>
    <h3>Response context</h3>
    <div>{{ String }}</div>
    <h2>Request context params</h2>
    <div>
        <ul>{{ #Context.params }}
        <li><pre>{{ . }}</pre></li>{{ /Context.params }}
        </ul>
    </div>
    <h2>Request context keys</h2>
    <div>
        <ul>{{ #Context.keys }}
        <li><pre>{{ . }}</pre></li>{{ /Context.keys }}
        </ul>
    </div>
    <h2>Request params</h2>
    <div>
        <ul>{{ #Params }}
        <li><pre>{{ . }}</pre></li>{{ /Params }}
        </ul>
    </div>
    <h2>Extra data</h2>
    <div>
        <ul>{{ #Extra }}
        <li><pre>{{ . }}</pre></li>{{ /Extra }}
        </ul>
    </div>
    <h2>Backend data</h2>
    <h3>Full response (as object)</h3>
    <div>
        <ul>{{ #Data }}
        <li><pre>{{ . }}</pre></li>{{ /Data }}
        </ul>
    </div>
    <h3>Full response (as array)</h3>
    <div>
        <ul>{{ #Array }}
        <li><pre>{{ . }}</pre></li>{{ /Array }}
        </ul>
    </div>
</div>`
)
