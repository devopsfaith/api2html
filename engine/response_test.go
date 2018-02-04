package engine

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestNoopResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	e := gin.New()
	e.GET("/", func(c *gin.Context) {
		resp, err := NoopResponse(c)
		if err != ErrNoResponseGeneratorDefined {
			t.Error("unexpected error:", err)
		}
		if len(resp.Data.Arr) != 0 {
			t.Error("unexpected response: %v", resp)
		}
		c.Status(200)
	})

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/", nil)
	e.ServeHTTP(w, r)
	if w.Result().StatusCode != 200 {
		t.Error("unexpected status code: %d", w.Result().StatusCode)
	}
}

func TestStaticResponseGenerator(t *testing.T) {
	subject := StaticResponseGenerator{Page{Extra: map[string]interface{}{"a": 42.0}}}
	gin.SetMode(gin.TestMode)
	e := gin.New()
	e.GET("/:first/:second", func(c *gin.Context) {
		resp, err := subject.ResponseGenerator(c)
		if err != nil {
			t.Error("unexpected error:", err.Error())
			return
		}
		checkCommonResponseProperties(t, resp)
		c.Status(200)
	})

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/foo/bar", nil)
	e.ServeHTTP(w, r)
	if w.Result().StatusCode != 200 {
		t.Error("unexpected status code: %d", w.Result().StatusCode)
	}
}

func TestDynamicResponseGenerator_koBackend(t *testing.T) {
	backendErr := fmt.Errorf("backendErr")
	expectedHeader := []string{"Header-Key", "header value"}
	subject := DynamicResponseGenerator{
		Page: Page{
			Extra:  map[string]interface{}{"a": 42.0},
			Header: expectedHeader[0],
		},
		Decoder: JSONDecoder,
		Backend: func(params map[string]string, headers map[string]string) (*http.Response, error) {
			if params["first"] != "foo" || params["second"] != "bar" {
				t.Error("unexpected params:", params)
			}
			if h, ok := headers[expectedHeader[0]]; !ok || h != expectedHeader[1] {
				t.Error("unexpected headers:", headers)
			}
			return nil, backendErr
		},
	}
	gin.SetMode(gin.TestMode)
	e := gin.New()
	e.GET("/:first/:second", func(c *gin.Context) {
		_, err := subject.ResponseGenerator(c)
		if err != backendErr {
			t.Error("unexpected error:", err)
			return
		}
		c.Status(200)
	})

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/foo/bar", nil)
	r.Header.Set(expectedHeader[0], expectedHeader[1])
	e.ServeHTTP(w, r)
	if w.Result().StatusCode != 200 {
		t.Error("unexpected status code: %d", w.Result().StatusCode)
	}
}

func TestDynamicResponseGenerator_koDecoder(t *testing.T) {
	decoderErr := fmt.Errorf("decoderErr")
	expectedResponse := "abcd"
	subject := DynamicResponseGenerator{
		Page: Page{Extra: map[string]interface{}{"a": 42.0}},
		Backend: func(_ map[string]string, _ map[string]string) (*http.Response, error) {
			return &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(expectedResponse))}, nil
		},
		Decoder: func(r io.Reader) (BackendData, error) {
			p := &bytes.Buffer{}
			p.ReadFrom(r)
			if p.String() != expectedResponse {
				t.Error("unexpected response:", p.String())
			}
			return BackendData{}, decoderErr
		},
	}
	gin.SetMode(gin.TestMode)
	e := gin.New()
	e.GET("/:first/:second", func(c *gin.Context) {
		_, err := subject.ResponseGenerator(c)
		if err != decoderErr {
			t.Error("unexpected error:", err)
			return
		}
		c.Status(200)
	})

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/foo/bar", nil)
	e.ServeHTTP(w, r)
	if w.Result().StatusCode != 200 {
		t.Error("unexpected status code: %d", w.Result().StatusCode)
	}
}

func TestDynamicResponseGenerator_ok(t *testing.T) {
	expectedResponse := "abcd"
	subject := DynamicResponseGenerator{
		Page: Page{Extra: map[string]interface{}{"a": 42.0}},
		Backend: func(_ map[string]string, _ map[string]string) (*http.Response, error) {
			return &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(expectedResponse))}, nil
		},
		Decoder: func(r io.Reader) (BackendData, error) {
			p := &bytes.Buffer{}
			p.ReadFrom(r)
			if p.String() != expectedResponse {
				t.Error("unexpected response:", p.String())
			}
			return BackendData{Obj: map[string]interface{}{"a": true}}, nil
		},
	}
	gin.SetMode(gin.TestMode)
	e := gin.New()
	e.GET("/:first/:second", func(c *gin.Context) {
		resp, err := subject.ResponseGenerator(c)
		if err != nil {
			t.Error("unexpected error:", err.Error())
			return
		}
		checkCommonResponseProperties(t, resp)

		if d, ok := resp.Data.Obj["a"].(bool); !ok || !d {
			t.Error("unexpected response. data: %v", resp.Data)
			return
		}
		c.Status(200)
	})

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/foo/bar", nil)
	e.ServeHTTP(w, r)
	if w.Result().StatusCode != 200 {
		t.Error("unexpected status code: %d", w.Result().StatusCode)
	}
}

func checkCommonResponseProperties(t *testing.T, resp ResponseContext) {
	if 42.0 != resp.Extra["a"].(float64) {
		t.Error("unexpected response. extra: %v", resp.Extra)
		return
	}
	if v, ok := resp.Params["first"]; !ok || v != "foo" {
		t.Error("unexpected response. first param: %v", resp.Params["first"])
		return
	}
	if v, ok := resp.Params["second"]; !ok || v != "bar" {
		t.Error("unexpected response. second param: %v", resp.Params["second"])
		return
	}
	if resp.Context == nil {
		t.Error("nil response context!")
		return
	}
	if v, ok := resp.Helper.(*tplHelper); !ok || v == nil {
		t.Error("unexpected response. helper: %v", resp.Helper)
		return
	}
}
