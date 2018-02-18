// Package engine contains all the required for building and running an API2HTML server
//
//	func Run(cfgPath string, devel bool) error {
//	 	errNilEngine := fmt.Errorf("serve cmd aborted: nil engine")
//	 	e, err := engine.New(cfgPath, devel)
//	 	if err != nil {
//	 		log.Println("engine creation aborted:", err.Error())
//	 		return err
//	 	}
//	 	if e == nil {
//	 		log.Println("engine creation aborted:", errNilEngine.Error())
//	 		return errNilEngine
//	 	}
//
//	 	time.Sleep(time.Second)
//
//	 	return eW.Run(fmt.Sprintf(":%d", port))
//	}
package engine

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Config is a struct with all the required definitions for building an API2HTML engine
type Config struct {
	Pages            []Page                 `json:"pages"`
	StaticTXTContent []string               `json:"static_txt_content"`
	Robots           bool                   `json:"robots"`
	Sitemap          bool                   `json:"sitemap"`
	Templates        map[string]string      `json:"templates"`
	Layouts          map[string]string      `json:"layouts"`
	Extra            map[string]interface{} `json:"extra"`
	PublicFolder     *PublicFolder          `json:"public_folder"`
}

// PublicFolder contains the info regarding the static contents to be served
type PublicFolder struct {
	Path   string `json:"path_to_folder"`
	Prefix string `json:"url_prefix"`
}

// Page defines the behaviour of the engine for a given URL pattern
type Page struct {
	Name              string
	URLPattern        string
	BackendURLPattern string
	Template          string
	Layout            string
	CacheTTL          string
	Header            string
	IsArray           bool
	Extra             map[string]interface{}
}

// New creates a gin engine with the default Factory
func New(cfgPath string, devel bool) (*gin.Engine, error) {
	return DefaultFactory.New(cfgPath, devel)
}

// Backend defines the signature of the function that creates a response for a request
// to a given backend
type Backend func(params map[string]string, headers map[string]string) (*http.Response, error)

// Renderer defines the interface for the template renderers
type Renderer interface {
	Render(io.Writer, interface{}) error
}

// RendererFunc is a function implementing the Renderer interface
type RendererFunc func(io.Writer, interface{}) error

// Render implements the Renderer interface
func (rf RendererFunc) Render(w io.Writer, v interface{}) error { return rf(w, v) }

// Subscription is a struct to be used to be notified after a change in the watched renderer
type Subscription struct {
	// Name is the name to watch
	Name string
	// In is the channel where the new renderer should be sent after a change
	In chan Renderer
}

// ErrorRenderer is a renderer that always returns the injected error
type ErrorRenderer struct {
	Error error
}

// Render implements the Renderer interface by returning the injected error
func (r ErrorRenderer) Render(_ io.Writer, _ interface{}) error { return r.Error }

// ErrNoResponseGeneratorDefined is the error returned when no ResponseGenerator has been defined
var ErrNoResponseGeneratorDefined = fmt.Errorf("no response generator defined")

// ErrNoBackendDefined is the error returned when no Backend has been defined
var ErrNoBackendDefined = fmt.Errorf("no backend defined")

// ErrNoRendererDefined is the error returned when no Renderer has been defined
var ErrNoRendererDefined = fmt.Errorf("no rendered defined")

// EmptyRenderer is the Renderer to be use if no other is defined
var EmptyRenderer = ErrorRenderer{ErrNoRendererDefined}
