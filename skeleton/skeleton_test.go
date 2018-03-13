package skeleton

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestNewSkel(t *testing.T) {
	defaultSkel := New("test")
	if err := defaultSkel.Create(); err != nil {
		t.Errorf("Error creating skeleton files: %s", err.Error())
	}
	defer os.RemoveAll("test")
	if err := defaultSkel.Create(); err.Error() != "output directory already exists" {
		t.Errorf("output directory already exists error should have been raised.")
	}
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
	hellodata, _ := ioutil.ReadFile("test/sources/global/static/hello.txt")
	expectedHello := string(`Hello, I am a text/plain content.`)
	if string(hellodata) != expectedHello {
		t.Error("Invalid content on hello.txt file.")
	}
}
