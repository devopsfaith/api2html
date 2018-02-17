package generator

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

func TestRender(t *testing.T) {
	dir, err := ioutil.TempDir("test", "render_output")
	if err != nil {
		log.Fatal(err)
	}

	defer os.RemoveAll(dir) // clean up

	data := Data{
		Config: map[string]Map{
			"cfg1": {
				"key1": "v1",
			},
		},
		I18N: map[string]Map{
			"I18N1": {
				"I18N1-1": "v1",
			},
		},
	}

	expectedErr := fmt.Errorf("some expected error")
	render := Render{
		OutputFolder: dir,
		Regexp:       nil,
		Dumper: func(_, _ string, _ Data) error {
			return expectedErr
		},
	}
	if err := render.Render("iso", data, dummyScanner{TmplFolder{Path: "path", Content: []string{"some_tmpl", "ignore_me"}}}); err != expectedErr {
		t.Error("render error:", err)
	}

	var total int
	render = Render{
		OutputFolder: dir,
		Regexp:       regexp.MustCompile("ignore"),
		Dumper: func(source, target string, d Data) error {
			if source != "path/some_tmpl" {
				t.Error("unexpected source:", source)
			}
			if !strings.Contains(target, "/iso/some_tmpl") || !strings.Contains(target, "test/render_output") {
				t.Error("unexpected target:", target)
			}
			if !reflect.DeepEqual(data, d) {
				t.Error("unexpected data:", d)
			}
			total++
			return nil
		},
	}

	if err := render.Render("iso", data, dummyScanner{TmplFolder{Path: "empty_path"}}); err != nil {
		t.Error("render error:", err)
	}

	if total != 0 {
		t.Errorf("unexpected number of calls to the dumper. have %d, want %d", total, 0)
	}

	if err := render.Render("iso", data, dummyScanner{TmplFolder{Path: "path", Content: []string{"some_tmpl", "ignore_me"}}}); err != nil {
		t.Error("render error:", err)
	}

	if total != 1 {
		t.Errorf("unexpected number of calls to the dumper. have %d, want %d", total, 1)
	}
}

// func TestRender(t *testing.T) {
// 	dir, err := ioutil.TempDir("test", "render_output")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	defer os.RemoveAll(dir) // clean up

// 	render := NewRender(dir, regexp.MustCompile("ignore"))

// 	tmpfn := filepath.Join(dir, "tmpfile")
// 	if err := ioutil.WriteFile(tmpfn, content, 0666); err != nil {
// 		log.Fatal(err)
// 	}
// }

type dummyScanner []TmplFolder

func (d dummyScanner) Scan() []TmplFolder {
	return d
}
