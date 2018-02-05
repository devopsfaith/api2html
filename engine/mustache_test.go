package engine

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestMustachePartials(t *testing.T) {
	fmt.Println(customPartialProvider.Get("api2html/debug"))
}

func TestNewMustacheRenderer(t *testing.T) {
	tmpl, err := NewMustacheRenderer(bytes.NewBufferString(`-{{ a }}-`))
	if err != nil {
		t.Error(err)
		return
	}

	if err := checkRenderer(tmpl); err != nil {
		t.Error(err)
	}
}

func TestNewLayoutMustacheRenderer(t *testing.T) {
	tmpl, err := NewLayoutMustacheRenderer(bytes.NewBufferString(`{{ a }}`), bytes.NewBufferString(`-{{{ content }}}-`))
	if err != nil {
		t.Error(err)
		return
	}

	if err := checkRenderer(tmpl); err != nil {
		t.Error(err)
	}
}

func TestNewMustacheRendererMap_ok(t *testing.T) {
	layoutPath := "a_layout.mustache"
	templatePath := "template.mustache"
	ioutil.WriteFile(layoutPath, []byte(`-{{{ content }}}-`), 0666)
	ioutil.WriteFile(templatePath, []byte(`-{{ a }}-`), 0666)
	renderers, err := NewMustacheRendererMap(Config{
		Templates: map[string]string{"t": templatePath},
		Layouts:   map[string]string{"l": layoutPath},
	})
	defer os.Remove(layoutPath)
	defer os.Remove(templatePath)
	if err != nil {
		t.Error(err)
		return
	}
	if _, ok := renderers["l"]; !ok {
		t.Error("layout renderer not found in the map")
	}
	tTmpl, ok := renderers["t"]
	if !ok {
		t.Error("template renderer not found in the map")
	}

	if err := checkRenderer(tTmpl); err != nil {
		t.Error(err)
	}
}

func TestNewMustacheRendererMap_koBadTemplate(t *testing.T) {
	layoutPath := "a_layout.mustache"
	ioutil.WriteFile(layoutPath, []byte(`-{{{ content`), 0666)
	_, err := NewMustacheRendererMap(Config{
		Layouts: map[string]string{"l": layoutPath},
	})
	defer os.Remove(layoutPath)
	if err == nil {
		t.Error("expecting error!")
		return
	}
}

func TestNewMustacheRendererMap_koNoFile(t *testing.T) {
	_, err := NewMustacheRendererMap(Config{
		Templates: map[string]string{"unknown": "unknown"},
	})
	if err == nil {
		t.Error("expecting error!")
		return
	}
}

func checkRenderer(tmpl Renderer) error {
	w := &bytes.Buffer{}
	ctx := map[string]interface{}{"a": 42}
	if err := tmpl.Render(w, ctx); err != nil {
		return err
	}
	if w.String() != "-42-" {
		return fmt.Errorf("unexpected render result: %s", w.String())
	}
	return nil
}
