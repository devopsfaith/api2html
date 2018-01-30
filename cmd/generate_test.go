package cmd

import (
	"fmt"
	"testing"

	"github.com/devopsfaith/api2html/generator"
)

func Test_defaultGeneratorFactory(t *testing.T) {
	g := defaultGeneratorFactory(".", "ignore")
	switch g.(type) {
	case *generator.BasicGenerator:
	default:
		t.Errorf("unexpected generator type: %t", g)
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
