package engine

import (
	"bytes"
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
	backend := DefaultClient(string(urlPattern))
	context, _ := gin.CreateTestContext(httptest.NewRecorder())
	_, err := backend(params, headers, context)
	if err != nil {
		t.Errorf("Backend response error: %s", err.Error())
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
