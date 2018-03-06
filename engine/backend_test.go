package engine

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

var (
	urlPattern = []byte("/test/:param")
	params     = map[string]string{
		"param": "replacetest",
	}
	headers = map[string]string{
		"X-Test": "testing",
	}
)

func TestNewBackend(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hi")
		if v, ok := r.Header["X-Test"]; ok {
			if v[0] != "testing" {
				t.Error("Invalid header content.")
			}
		}
		if r.URL.RequestURI() != "/test/replacetest" {
			fmt.Println(r.URL.RequestURI())
			t.Error("Invalid URL.")
		}
	}))
	defer mockServer.Close()
	backend := DefaultClient(fmt.Sprintf("%s%s", mockServer.URL, string(urlPattern)))
	context, _ := gin.CreateTestContext(httptest.NewRecorder())
	resp, err := backend(params, headers, context)
	if err != nil {
		t.Errorf("Backend response error: %s", err.Error())
	}
	if resp.StatusCode != 200 {
		t.Error("Invalid status code.")
	}
}

func TestReplaceParams(t *testing.T) {

	expectedResult := []byte("/test/replacetest")
	// Test empty params
	if !bytes.Equal(urlPattern, replaceParams(urlPattern, map[string]string{})) {
		t.Error("An empty param list should return the same URLPattern")
	}

	// Test replace params
	if !bytes.Equal(expectedResult, replaceParams(urlPattern, params)) {
		t.Error("The replace is not working as expected.")
	}
}
