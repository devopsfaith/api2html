package engine

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestFactory_New_koConfigParser(t *testing.T) {
	expectedErr := fmt.Errorf("boooom")
	ef := Factory{
		Parser: func(path string) (Config, error) {
			if path != "something" {
				t.Errorf("unexpected path: %s", path)
			}
			return Config{}, expectedErr
		},
	}
	if _, err := ef.New("something", true); err == nil {
		t.Error("expecting error")
	} else if err != expectedErr {
		t.Errorf("unexpected error: %s", err.Error())
	}
}

func TestFactory_New_ok(t *testing.T) {
	if err := ioutil.WriteFile("test_tmpl", []byte("hi, {{Extra.name}}!"), 0644); err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if err := ioutil.WriteFile("test_lyt", []byte("-{{{content}}}-"), 0644); err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	defer os.Remove("test_tmpl")
	defer os.Remove("test_lyt")
	expectedCfg := Config{
		Pages: []Page{
			{
				URLPattern: "/a",
				Layout:     "b",
				Template:   "a",
				Extra: map[string]interface{}{
					"name": "stranger",
				},
			},
		},
		Templates: map[string]string{"a": "test_tmpl"},
		Layouts:   map[string]string{"b": "test_lyt"},
	}
	templateStore := NewTemplateStore()
	ef := DefaultFactory
	ef.Parser = func(path string) (Config, error) {
		if path != "something" {
			t.Errorf("unexpected path: %s", path)
		}
		return expectedCfg, nil
	}
	ef.TemplateStoreFactory = func() *TemplateStore { return templateStore }
	ef.MustachePageFactory = func(e *gin.Engine, ts *TemplateStore) MustachePageFactory {
		if ts != templateStore {
			t.Errorf("unexpected template store: %v", ts)
		}
		return NewMustachePageFactory(e, ts)
	}

	e, err := ef.New("something", true)
	if err != nil {
		t.Errorf("unexpected error: %s", err.Error())
		return
	}

	time.Sleep(200 * time.Millisecond)

	assertResponse(t, e, "/a", http.StatusOK, "-hi, stranger!-")
	assertResponse(t, e, "/b", http.StatusNotFound, default404Tmpl)
}
