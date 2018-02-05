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

func ExampleResponseContext_String() {
	r := ResponseContext{
		Data: map[string]interface{}{
			"a": "foo",
			"b": 42,
		},
		Params: map[string]string{"p1": "v1"},
		Extra: map[string]interface{}{
			"extra1": "foo",
			"extra2": 42,
		},
	}
	fmt.Println(r.String())
	// Output:
	// {
	// 	"Data": {
	// 		"a": "foo",
	// 		"b": 42
	// 	},
	// 	"Array": null,
	// 	"Extra": {
	// 		"extra1": "foo",
	// 		"extra2": 42
	// 	},
	// 	"Params": {
	// 		"p1": "v1"
	// 	}
	// }
}

func TestNoopResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	e := gin.New()
	e.GET("/", func(c *gin.Context) {
		resp, err := NoopResponse(c)
		if err != ErrNoResponseGeneratorDefined {
			t.Error("unexpected error:", err)
		}
		if len(resp.Array) != 0 {
			t.Errorf("unexpected response: %v", resp)
		}
		if len(resp.Data) != 0 {
			t.Errorf("unexpected response: %v", resp)
		}
		c.Status(200)
	})

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/", nil)
	e.ServeHTTP(w, r)
	if w.Result().StatusCode != 200 {
		t.Errorf("unexpected status code: %d", w.Result().StatusCode)
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
		t.Errorf("unexpected status code: %d", w.Result().StatusCode)
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
		t.Errorf("unexpected status code: %d", w.Result().StatusCode)
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
		Decoder: func(r io.Reader, c *ResponseContext) error {
			p := &bytes.Buffer{}
			p.ReadFrom(r)
			if p.String() != expectedResponse {
				t.Error("unexpected response:", p.String())
			}
			return decoderErr
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
		t.Errorf("unexpected status code: %d", w.Result().StatusCode)
	}
}

func TestDynamicResponseGenerator_ok(t *testing.T) {
	expectedResponse := "abcd"
	subject := DynamicResponseGenerator{
		Page: Page{Extra: map[string]interface{}{"a": 42.0}},
		Backend: func(_ map[string]string, _ map[string]string) (*http.Response, error) {
			return &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(expectedResponse))}, nil
		},
		Decoder: func(r io.Reader, c *ResponseContext) error {
			p := &bytes.Buffer{}
			p.ReadFrom(r)
			if p.String() != expectedResponse {
				t.Error("unexpected response:", p.String())
			}
			c.Data = map[string]interface{}{"a": true}
			return nil
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

		if d, ok := resp.Data["a"].(bool); !ok || !d {
			t.Errorf("unexpected response. data: %v", resp.Data)
			return
		}
		c.Status(200)
	})

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/foo/bar", nil)
	e.ServeHTTP(w, r)
	if w.Result().StatusCode != 200 {
		t.Errorf("unexpected status code: %d", w.Result().StatusCode)
	}
}

func checkCommonResponseProperties(t *testing.T, resp ResponseContext) {
	if 42.0 != resp.Extra["a"].(float64) {
		t.Errorf("unexpected response. extra: %v", resp.Extra)
		return
	}
	if v, ok := resp.Params["first"]; !ok || v != "foo" {
		t.Errorf("unexpected response. first param: %v", resp.Params["first"])
		return
	}
	if v, ok := resp.Params["second"]; !ok || v != "bar" {
		t.Errorf("unexpected response. second param: %v", resp.Params["second"])
		return
	}
	if resp.Context == nil {
		t.Error("nil response context!")
		return
	}
	if v, ok := resp.Helper.(*tplHelper); !ok || v == nil {
		t.Errorf("unexpected response. helper: %v", resp.Helper)
		return
	}
}
