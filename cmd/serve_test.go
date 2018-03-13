package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/devopsfaith/api2html/engine"
	"github.com/gin-gonic/gin"
)

func Test_defaultEngineFactory(t *testing.T) {
	cfg := engine.Config{}
	f, err := ioutil.TempFile(".", "")
	if err != nil {
		t.Error(err)
		return
	}
	if jerr := json.NewEncoder(f).Encode(&cfg); jerr != nil {
		t.Error(jerr)
		return
	}
	f.Close()

	defer os.Remove(f.Name())

	g, err := defaultEngineFactory(f.Name(), true)
	if err != nil {
		t.Errorf("getting the default engine: %s", err.Error())
		return
	}
	switch g.(type) {
	case *gin.Engine:
	default:
		t.Errorf("unexpected engine type: %T", g)
	}
}

func Test_serveWrapper_koErroredEngineFactory(t *testing.T) {
	expectedError := fmt.Errorf("expect me!")
	subject := serveWrapper{erroredEngineFactory(expectedError)}

	if err := subject.Serve(nil, []string{}); err == nil {
		t.Error("expecting error!")
		return
	} else if err != expectedError {
		t.Errorf("unexpected error! want: %s, got: %s", expectedError.Error(), err.Error())
		return
	}
}

func Test_serveWrapper_koErroredEngine(t *testing.T) {
	subject := serveWrapper{customEngineFactory(nil)}

	if err := subject.Serve(nil, []string{}); err == nil {
		t.Error("expecting error!")
		return
	} else if err != errNilEngine {
		t.Errorf("unexpected error! want: %s, got: %s", errNilEngine.Error(), err.Error())
		return
	}

	expectedError := fmt.Errorf("expect me!")
	subject = serveWrapper{customEngineFactory(erroredEngine{expectedError})}

	if err := subject.Serve(nil, []string{}); err == nil {
		t.Error("expecting error!")
	} else if err != expectedError {
		t.Errorf("unexpected error! want: %s, got: %s", expectedError.Error(), err.Error())
	}
}

func erroredEngineFactory(err error) engineFactory {
	return func(_ string, _ bool) (engineWrapper, error) { return nil, err }
}

func customEngineFactory(e engineWrapper) engineFactory {
	return func(_ string, _ bool) (engineWrapper, error) { return e, nil }
}

type erroredEngine struct {
	err error
}

func (e erroredEngine) Run(_ ...string) error {
	return e.err
}
