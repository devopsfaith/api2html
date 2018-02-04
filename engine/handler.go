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
	result, err := h.ResponseGenerator(c)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.Header("Cache-Control", h.CacheControl)
	if err := h.Renderer.Render(c.Writer, &result); err != nil {
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
