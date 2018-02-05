package engine

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

func TestNewStaticHandler(t *testing.T) {
	fileName := fmt.Sprintf("testErrorHAndler-%d", time.Now().Unix())
	data := []byte("sample data to be dumped by the error handler")
	err := ioutil.WriteFile(fileName, data, 0666)
	if err != nil {
		t.Error(err)
		return
	}
	defer os.Remove(fileName)

	eh, err := NewStaticHandler(fileName)
	if err != nil {
		t.Error(err)
		return
	}

	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.GET("/static", eh.HandlerFunc())

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/static", nil)
	if err != nil {
		t.Error(err)
		return
	}
	engine.ServeHTTP(w, req)

	if w.Result().StatusCode != 200 {
		t.Errorf("unexpected status code: %d", w.Result().StatusCode)
	}
	res, err := ioutil.ReadAll(w.Result().Body)
	if err != nil {
		t.Error(err)
		return
	}
	w.Result().Body.Close()
	if string(res) != string(data) {
		t.Errorf("unexpected response content: %s", string(res))
	}
}

func TestNewErrorHandler(t *testing.T) {
	fileName := fmt.Sprintf("testErrorHAndler-%d", time.Now().Unix())
	data := []byte("sample data to be dumped by the error handler")
	err := ioutil.WriteFile(fileName, data, 0666)
	if err != nil {
		t.Error(err)
		return
	}
	defer os.Remove(fileName)

	eh, err := NewErrorHandler(fileName)
	if err != nil {
		t.Error(err)
		return
	}

	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.GET("/middleware/ok", eh.HandlerFunc(), func(c *gin.Context) { c.String(200, "hi there!") })
	engine.GET("/middleware/ko", eh.HandlerFunc(), func(c *gin.Context) { c.AbortWithStatus(987) })

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/middleware/ok", nil)
	if err != nil {
		t.Error(err)
		return
	}
	engine.ServeHTTP(w, req)

	if w.Result().StatusCode != 200 {
		t.Errorf("unexpected status code: %d", w.Result().StatusCode)
	}
	res, err := ioutil.ReadAll(w.Result().Body)
	if err != nil {
		t.Error(err)
		return
	}
	w.Result().Body.Close()
	if string(res) != "hi there!" {
		t.Errorf("unexpected response content: %s", string(res))
	}

	w = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/middleware/ko", nil)
	if err != nil {
		t.Error(err)
		return
	}
	engine.ServeHTTP(w, req)

	if w.Result().StatusCode != 987 {
		t.Errorf("unexpected status code: %d", w.Result().StatusCode)
	}
	res, err = ioutil.ReadAll(w.Result().Body)
	if err != nil {
		t.Error(err)
		return
	}
	w.Result().Body.Close()
	if string(res) != string(data) {
		t.Errorf("unexpected response content: %s", string(res))
	}
}
