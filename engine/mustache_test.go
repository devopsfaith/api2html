package engine

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
	"testing/iotest"
)

func TestMustachePartials(t *testing.T) {
	tmpl, err := customPartialProvider.Get("api2html/debug")
	if err != nil {
		t.Error(err)
	}
	if tmpl == "" {
		t.Error("empty partial")
	}
	tmpl, err = customPartialProvider.Get("_____unnknown_____")
	if err != nil {
		t.Error(err)
	}
	if tmpl != "" {
		t.Error("unexpected result:", tmpl)
	}
}

func TestNewMustacheRenderer_ok(t *testing.T) {
	tmpl, err := NewMustacheRenderer(bytes.NewBufferString(`-{{ a }}-`))
	if err != nil {
		t.Error(err)
		return
	}

	if err := checkRenderer(tmpl); err != nil {
		t.Error(err)
	}
}

func TestNewMustacheRenderer_ko(t *testing.T) {
	_, err := NewMustacheRenderer(bytes.NewBufferString(`-{{ a `))
	if err == nil {
		t.Error("expecting error")
	}
}

func TestNewLayoutMustacheRenderer_ok(t *testing.T) {
	tmpl, err := NewLayoutMustacheRenderer(bytes.NewBufferString(`{{ a }}`), bytes.NewBufferString(`-{{{ content }}}-`))
	if err != nil {
		t.Error(err)
		return
	}

	if err := checkRenderer(tmpl); err != nil {
		t.Error(err)
	}
}

func TestNewLayoutMustacheRenderer_ko(t *testing.T) {
	_, err := NewLayoutMustacheRenderer(bytes.NewBufferString(`{{ a `), bytes.NewBufferString(`-{{{ content }}}-`))
	if err == nil {
		t.Error("expecting error")
	}
	_, err = NewLayoutMustacheRenderer(bytes.NewBufferString(`{{ a }}`), bytes.NewBufferString(`-{{{ content -`))
	if err == nil {
		t.Error("expecting error")
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

func Test_newMustacheTemplate(t *testing.T) {
	b := make([]byte, 1024)
	rand.Read(b)
	if _, err := newMustacheTemplate(iotest.TimeoutReader(bytes.NewBuffer(b))); err == nil {
		t.Error("expecting error!")
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
