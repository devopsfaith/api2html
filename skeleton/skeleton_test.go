package skeleton

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestNewBlogSkel(t *testing.T) {
	defaultSkel := NewBlog("test")
	if err := defaultSkel.Create(); err != nil {
		t.Errorf("Error creating skeleton files: %s", err.Error())
	}
	defer os.RemoveAll("test")

	counter := 0
	err := filepath.Walk("test", func(path string, f os.FileInfo, err error) error {
		if _, err := os.Stat("test"); err != nil {
			return err
		}
		if !f.IsDir() {
			counter++
		}
		return nil
	})
	if err != nil {
		t.Errorf("Problem walking the test directory: %s", err.Error())
	}
	if counter != 19 {
		t.Error("File count wrong, the test should have been generated 19 files.")
	}
	hellodata, _ := ioutil.ReadFile("test/blog/sources/global/static/hello.txt")
	expectedHello := string(`Hello, I am a text/plain content.`)
	if string(hellodata) != expectedHello {
		t.Error("Invalid content on hello.txt file.")
	}
}

func TestNewSkel(t *testing.T) {
	defaultSkel := New("test", []string{"/blog/i18n/es_ES.ini", "/blog/sources/global/static/hello.txt"})
	if err := defaultSkel.Create(); err != nil {
		t.Errorf("Error creating skeleton files: %s", err.Error())
	}
	defer os.RemoveAll("test")

	counter := 0
	err := filepath.Walk("test", func(path string, f os.FileInfo, err error) error {
		if _, err := os.Stat("test"); err != nil {
			return err
		}
		if !f.IsDir() {
			counter++
		}
		return nil
	})
	if err != nil {
		t.Errorf("Problem walking the test directory: %s", err.Error())
	}
	if counter != 2 {
		t.Error("File count wrong, the test should have been generated 19 files.")
	}
	hellodata, _ := ioutil.ReadFile("test/blog/sources/global/static/hello.txt")
	expectedHello := string(`Hello, I am a text/plain content.`)
	if string(hellodata) != expectedHello {
		t.Error("Invalid content on hello.txt file.")
	}
}
