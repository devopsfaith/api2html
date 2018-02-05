package engine

import (
	"fmt"
	"io"
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

func TestNewStaticHandler_ko(t *testing.T) {
	_, err := NewStaticHandler("unknown_file_not_present_in_the_fs")
	if err == nil {
		t.Error("error expected")
	}
}

func TestNewErrorHandlerr_ko(t *testing.T) {
	_, err := NewErrorHandler("unknown_file_not_present_in_the_fs")
	if err == nil {
		t.Error("error expected")
	}
}

func TestNewHandler(t *testing.T) {
	responseCtx := ResponseContext{
		Array: []map[string]interface{}{
			map[string]interface{}{"a": "foo"},
		},
	}
	layout := "layout"
	templateName := "name"
	responseBody := "some response content"
	cfg := HandlerConfig{
		Renderer: EmptyRenderer,
		ResponseGenerator: func(_ *gin.Context) (ResponseContext, error) {
			return responseCtx, nil
		},
		Page: Page{
			Template: templateName,
			Layout:   layout,
		},
	}
	subscriptionChan := make(chan Subscription)
	h := NewHandler(cfg, subscriptionChan)

	subscription := <-subscriptionChan
	if subscription.Name != layout+"-:-"+templateName {
		t.Errorf("unexpected subscription topic: %s", subscription.Name)
		return
	}
	subscription.In <- RendererFunc(func(w io.Writer, v interface{}) error {
		if tmp, ok := v.(ResponseContext); !ok {
			t.Errorf("unexpected type %t", v)
			return nil
		} else if len(tmp.Array) != 1 {
			t.Errorf("unexpected value %v", tmp)
			return nil
		}
		_, err := w.Write([]byte(responseBody))
		return err
	})
	<-subscriptionChan

	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.GET("/", h.HandlerFunc)

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
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
	if string(res) != responseBody {
		t.Errorf("unexpected response content: %s", string(res))
	}
}
