package engine

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	api2htmlhttp "github.com/devopsfaith/api2html/engine/http"
	"github.com/gin-gonic/gin"
)

// HandlerConfig defines a Handler
type HandlerConfig struct {
	// Page contains the page description
	Page Page
	// Renderer is the component responsible for rendering the responses
	Renderer Renderer
	// ResponseGenerator gets the data required for generating a response
	// it can get it from a static, local source or from a remote api
	// endpoint
	ResponseGenerator ResponseGenerator
	// CacheControl is the Cache-Control string added into the response headers
	// if everything goes ok
	CacheControl string
}

// DefaultHandlerConfig contains the dafult values for a HandlerConfig
var DefaultHandlerConfig = HandlerConfig{
	Page{},
	EmptyRenderer,
	NoopResponse,
	"public, max-age=3600",
}

// ResponseGenerator is a function that, given a gin request, returns a response struc and an error
//
// The returned response, even being a map[string]interface{}, can contain these fields by convention:
// 	- extra:	The content of the 'extra' property of the page
// 	- context:	The gin context for the request
// 	- params:	The params of the request
// 	- helper:	A struct containing a few basic template helpers
// 	- data:		Depending on the response generator implementation, it cotains the backend data
type ResponseGenerator func(*gin.Context) (map[string]interface{}, error)

// NoopResponse is a ResponseGenerator that always returns an empty response and the
// ErrNoResponseGeneratorDefined error
func NoopResponse(_ *gin.Context) (map[string]interface{}, error) {
	return map[string]interface{}{}, ErrNoResponseGeneratorDefined
}

// StaticResponseGenerator is a ResponseGenerator that creates a response just by adding the
// default response values
// 	map[string]interface{}{
// 		"extra":   s.Page.Extra,
// 		"context": c,
// 		"params":  params,
// 		"helper":  &tplHelper{},
// 	}
type StaticResponseGenerator struct {
	Page Page
}

// ResponseGenerator implements the ResponseGenerator interface
func (s *StaticResponseGenerator) ResponseGenerator(c *gin.Context) (map[string]interface{}, error) {
	params := map[string]string{}
	for _, v := range c.Params {
		params[v.Key] = v.Value
	}
	target := map[string]interface{}{
		"extra":   s.Page.Extra,
		"context": c,
		"params":  params,
		"helper":  &tplHelper{},
	}
	return target, nil
}

// DynamicResponseGenerator is a ResponseGenerator that creates a response by adding the decoded data
// returned by the Backend wo the default response values. Depending on the selected decoder,
// the generated responses may have this structure
// 	map[string]interface{}{
// 		"data":    []map[sitring]interface{}{},
// 		"extra":   s.Page.Extra,
// 		"context": c,
// 		"params":  params,
// 		"helper":  &tplHelper{},
// 	}
// or this one
// 	map[string]interface{}{
// 		"data":    []map[sitring]interface{}{},
// 		"extra":   s.Page.Extra,
// 		"context": c,
// 		"params":  params,
// 		"helper":  &tplHelper{},
// 	}
type DynamicResponseGenerator struct {
	Page    Page
	Backend Backend
	Decoder Decoder
}

// ResponseGenerator implements the ResponseGenerator interface
func (drg *DynamicResponseGenerator) ResponseGenerator(c *gin.Context) (map[string]interface{}, error) {
	params := map[string]string{}
	for _, v := range c.Params {
		params[v.Key] = v.Value
	}
	headers := map[string]string{}
	h := c.Request.Header.Get(drg.Page.Header)
	if h != "" {
		headers[drg.Page.Header] = h
	}
	resp, err := drg.Backend(params, headers)
	if err != nil {
		return map[string]interface{}{}, err
	}
	defer resp.Body.Close()

	target, err := drg.Decoder(resp.Body)
	if err != nil {
		return map[string]interface{}{}, err
	}
	target["extra"] = drg.Page.Extra
	target["context"] = c
	target["params"] = params
	target["helper"] = &tplHelper{}
	return target, nil
}

// NewHandlerConfig creates a HandlerConfig from the given Page definition
func NewHandlerConfig(page Page) HandlerConfig {
	d, err := time.ParseDuration(page.CacheTTL)
	if err != nil {
		d = time.Hour
	}
	cacheTTL := fmt.Sprintf("public, max-age=%d", int(d.Seconds()))

	if page.BackendURLPattern == "" {
		rg := StaticResponseGenerator{page}
		return HandlerConfig{
			page,
			DefaultHandlerConfig.Renderer,
			rg.ResponseGenerator,
			cacheTTL,
		}
	}

	decoder := JSONDecoder
	if page.IsArray {
		decoder = JSONArrayDecoder
	}
	rg := DynamicResponseGenerator{page, Backend(api2htmlhttp.CachedClient(page.BackendURLPattern)), decoder}

	return HandlerConfig{
		page,
		DefaultHandlerConfig.Renderer,
		rg.ResponseGenerator,
		cacheTTL,
	}
}

// NewHandler creates a Handler with the given configuration. The returned handler will be keeping itself
// subscribed to the latest template updates using the given subscription channel, allowing hot
// template reloads
func NewHandler(cfg HandlerConfig, subscriptionChan chan Subscription) *Handler {
	h := &Handler{
		cfg.Page,
		cfg.Renderer,
		make(chan Renderer),
		subscriptionChan,
		cfg.ResponseGenerator,
		cfg.CacheControl,
	}
	go h.updateRenderer()
	return h
}

// Handler is a struct that combines a renderer and a response generator for handling
// http requests.
//
// The handler is able to keep itself subscribed to the last renderer version to use
// by wrapping its Input channel into a Subscription and sending it throught the Subscribe
// channel every time it gets a new Renderer
type Handler struct {
	Page              Page
	Renderer          Renderer
	Input             chan Renderer
	Subscribe         chan Subscription
	ResponseGenerator ResponseGenerator
	CacheControl      string
}

func (h *Handler) updateRenderer() {
	topic := h.Page.Template
	if h.Page.Layout != "" {
		topic = fmt.Sprintf("%s-:-%s", h.Page.Layout, h.Page.Template)
	}
	for {
		h.Subscribe <- Subscription{topic, h.Input}
		h.Renderer = <-h.Input
	}
}

// HandlerFunc handles a gin request rendering the data returned by the response generator.
// If the response generator does not return an error, it adds a Cache-Control header
func (h *Handler) HandlerFunc(c *gin.Context) {
	target, err := h.ResponseGenerator(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Header("Cache-Control", h.CacheControl)
	if err := h.Renderer.Render(c.Writer, &target); err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
}

// NewErrorHandler creates a ErrorHandler using the content of the received path
func NewErrorHandler(path string) (ErrorHandler, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println("reading", path, ":", err.Error())
		return ErrorHandler{}, err
	}
	return ErrorHandler{data}, nil
}

// ErrorHandler is a Handler that writes the injected content. It's intended to be dispatched
// by the gin special handlers (NoRoute, NoMethod) but they can also be used as regular handlers
type ErrorHandler struct {
	Content []byte
}

// HandlerFunc is a gin middleware for dealing with some errors
func (e *ErrorHandler) HandlerFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if !c.IsAborted() {
			return
		}

		c.Writer.Write(e.Content)
	}
}

// StaticHandlerFunc creates a gin handler that does nothing but writing the static content
func (e *ErrorHandler) StaticHandlerFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Write(e.Content)
	}
}

type tplHelper struct {
}

func (tplHelper) Now() string {
	return time.Now().String()
}
