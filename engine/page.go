package engine

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// NewMustachePageFactory creates a MustachePageFactory with the injected params
func NewMustachePageFactory(e *gin.Engine, ts *TemplateStore) MustachePageFactory {
	return MustachePageFactory{e, ts}
}

// MustachePageFactory is a component that sets up the gin engine and the template store
type MustachePageFactory struct {
	Engine        *gin.Engine
	TemplateStore *TemplateStore
}

// Build sets up the injected gin engine and template store depending on the contents of
// the received configuration
func (m *MustachePageFactory) Build(cfg Config) {
	templates, err := NewMustacheRendererMap(cfg)
	if err != nil {
		panic(err)
	}

	for _, page := range cfg.Pages {
		h := NewHandler(NewHandlerConfig(page), m.TemplateStore.Subscribe)
		m.Engine.GET(page.URLPattern, h.HandlerFunc)

		time.Sleep(100 * time.Millisecond)

		r, ok := templates[page.Template]
		if !ok {
			fmt.Println("handler without template", page.Name, page.Template)
			continue
		}
		m.TemplateStore.Set(page.Template, r)
		if page.Layout == "" {
			fmt.Println("handler without layout", page.Name, page.Layout)
			continue
		}
		l, ok := templates[page.Layout]
		if !ok {
			fmt.Println("layout not defined", page.Layout)
			continue
		}
		m.TemplateStore.Set(page.Layout, l)

		m.TemplateStore.Set(fmt.Sprintf("%s-:-%s", h.Page.Layout, h.Page.Template), &LayoutMustacheRenderer{r.tmpl, l.tmpl})
	}
}
