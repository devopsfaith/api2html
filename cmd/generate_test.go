package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/devopsfaith/api2html/generator"
)

func Test_defaultGeneratorFactory(t *testing.T) {
	g := defaultGeneratorFactory(".", "ignore")
	switch g.(type) {
	case *generator.BasicGenerator:
	default:
		t.Errorf("unexpected generator type: %T", g)
	}
}

func Test_generatorWrapper_koErroredGenerator(t *testing.T) {
	expectedError := fmt.Errorf("expect me!")

	subject := generatorWrapper{func(_, _ string) generator.Generator {
		return erroredGenerator{expectedError}
	}}

	if err := subject.Generate(nil, []string{}); err == nil {
		t.Error("expecting error!")
	} else if err != expectedError {
		t.Errorf("unexpected error! want: %s, got: %s", expectedError.Error(), err.Error())
	}
}

func Test_generatorWrapper(t *testing.T) {
	expectedIso := "iso1"
	expectedBasePath := "basepath1"
	expectedIgnoreRegex := "ignoreRegex1"

	wrongIso := "no_iso"
	wrongBasePath := "no_base_path"
	wrongIgnoreRegex := "no_ignore_regex"

	errIso := fmt.Errorf("generating. have: %s, want: %s", wrongIso, expectedIso)
	errBasePath := fmt.Errorf("wrong base path. have: %s, want: %s", wrongBasePath, expectedBasePath)
	errIgnoreRegex := fmt.Errorf("wrong ignore regex. have: %s, want: %s", wrongIgnoreRegex, expectedIgnoreRegex)

	spy := func(bp, reg string) generator.Generator {
		if expectedBasePath != bp {
			return erroredGenerator{fmt.Errorf("wrong base path. have: %s, want: %s", bp, expectedBasePath)}
		}
		if expectedIgnoreRegex != reg {
			return erroredGenerator{fmt.Errorf("wrong ignore regex. have: %s, want: %s", reg, expectedIgnoreRegex)}
		}
		return spyGenerator{expectedIso}
	}

	subject := generatorWrapper{spy}

	isos = wrongIso
	basePath = wrongBasePath
	ignoreRegex = wrongIgnoreRegex

	if err := subject.Generate(nil, []string{}); err == nil {
		t.Error("expecting error!")
		return
	} else if err.Error() != errBasePath.Error() {
		t.Errorf("unexpected error! want: %s, got: %s", errBasePath.Error(), err.Error())
		return
	}

	basePath = expectedBasePath

	if err := subject.Generate(nil, []string{}); err == nil {
		t.Error("expecting error!")
		return
	} else if err.Error() != errIgnoreRegex.Error() {
		t.Errorf("unexpected error! want: %s, got: %s", errIgnoreRegex.Error(), err.Error())
		return
	}

	ignoreRegex = expectedIgnoreRegex

	if err := subject.Generate(nil, []string{}); err == nil {
		t.Error("expecting error!")
		return
	} else if err.Error() != errIso.Error() {
		t.Errorf("unexpected error! want: %s, got: %s", errIso.Error(), err.Error())
		return
	}

	isos = expectedIso

	if err := subject.Generate(nil, []string{}); err != nil {
		t.Error("unexpected error:", err.Error())
	}
}

func Test_generatorWatchWrapper_koErroredGenerator(t *testing.T) {
	expectedError := fmt.Errorf("expect me!")

	subject := generatorWatchWrapper{generatorWrapper{func(_, _ string) generator.Generator {
		return erroredGenerator{expectedError}
	}}}

	if err := subject.Watch(nil, []string{}); err == nil {
		t.Error("expecting error!")
	} else if err != expectedError {
		t.Errorf("unexpected error! want: %s, got: %s", expectedError.Error(), err.Error())
	}
}

func Test_generatorWatchWrapper_koErroredGeneratorAfterChange(t *testing.T) {
	name, err := ioutil.TempDir(".", "tmp")
	if err != nil {
		t.Error(err)
		return
	}
	defer os.RemoveAll(name)

	expectedError := fmt.Errorf("expect me!")
	var counter uint64
	isos = "*"
	basePath = name
	ignoreRegex = "ignore"

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	subject := generatorWatchWrapper{generatorWrapper{func(basePath, ignoreRegex string) generator.Generator {
		select {
		case <-ctx.Done():
			return erroredGenerator{expectedError}
		default:
		}

		if basePath != name {
			t.Errorf("unexpected base path. have %s want %s", basePath, name)
		}
		if ignoreRegex != "ignore" {
			t.Errorf("unexpected ignore regex. have %s want %s", ignoreRegex, "ignore")
		}
		atomic.AddUint64(&counter, 1)
		fmt.Println("generation triggered!", counter)
		if atomic.LoadUint64(&counter) < 2 {
			return spyGenerator{isos}
		}
		return erroredGenerator{expectedError}
	}}}

	wg := &sync.WaitGroup{}
	go func() {
		wg.Add(1)
		if werr := subject.Watch(nil, []string{}); werr == nil {
			t.Error("expecting error!")
		} else if werr != expectedError {
			t.Errorf("unexpected error! want: %s, got: %s", expectedError.Error(), werr.Error())
		}
		wg.Done()
	}()

	time.Sleep(150 * time.Millisecond)

	if atomic.LoadUint64(&counter) != 1 {
		t.Errorf("unexpected number of calls to the genetator. have %d, want %d", atomic.LoadUint64(&counter), 1)
	}
	if err = ioutil.WriteFile(name+"/test", []byte("12345678"), 0644); err != nil {
		t.Error(err)
		return
	}

	time.Sleep(150 * time.Millisecond)

	if atomic.LoadUint64(&counter) != 2 {
		t.Errorf("unexpected number of calls to the genetator. have %d, want %d", atomic.LoadUint64(&counter), 2)
	}

	cancel()
}

type erroredGenerator struct {
	err error
}

func (e erroredGenerator) Generate(_ string) error {
	return e.err
}

type spyGenerator struct {
	want string
}

func (s spyGenerator) Generate(have string) error {
	if s.want != have {
		return fmt.Errorf("generating. have: %s, want: %s", have, s.want)
	}
	return nil
}
