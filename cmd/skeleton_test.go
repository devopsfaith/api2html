package cmd

import (
	"fmt"
	"os"
	"testing"

	"github.com/devopsfaith/api2html/skeleton"
)

func Test_defaultSkelFactory(t *testing.T) {
	g := defaultBlogSkelFactory("test")
	switch g.(type) {
	case skeleton.Skel:
	default:
		t.Errorf("unexpected generator type: %T", g)
	}
}

func Test_skelWrapper_koErroredSkel(t *testing.T) {
	expectedError := fmt.Errorf("expect me")

	skel := skelWrapper{func(_ string) skeleton.Skel {
		return erroredSkel{expectedError}
	}}

	if err := skel.Create(nil, []string{}); err == nil {
		t.Error("expecting error!")
	} else if err != expectedError {
		t.Errorf("unexpected error! want: %s, got: %s", expectedError.Error(), err.Error())
	}
}

func Test_skelWrapper(t *testing.T) {

	skel := skelWrapper{func(_ string) skeleton.Skel {
		return simpleSkel{outputPath}
	}}

	defer os.Remove("example")
	if err := skel.Create(nil, []string{}); err != nil {
		t.Errorf("unexpected error: %s", err.Error())
	}
	if _, err := os.Stat("example"); err != nil {
		t.Errorf("cannot locate test output dir: %s", err.Error())
	}
}

type erroredSkel struct {
	err error
}

func (e erroredSkel) Create() error {
	return e.err
}

type simpleSkel struct {
	outputPath string
}

func (s simpleSkel) Create() error {
	return os.Mkdir(s.outputPath, os.ModePerm)
}
