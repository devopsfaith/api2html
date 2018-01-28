package engine

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	api2htmlhttp "github.com/devopsfaith/api2html/engine/http"
	"github.com/gin-gonic/gin"
)

type HandlerConfig struct {
	Page              Page
	Renderer          Renderer
	ResponseGenerator ResponseGenerator
	CacheControl      string
}

var DefaultHandlerConfig = HandlerConfig{
	Page{},
	EmptyRenderer,
	NoopResponse,
	"public, max-age=3600",
}

type ResponseGenerator func(*gin.Context) (map[string]interface{}, error)

func NoopResponse(_ *gin.Context) (map[string]interface{}, error) {
	return map[string]interface{}{}, ErrNoResponseGeneratorDefined
}

type StaticResponseGenerator struct {
	Page Page
}

func (s *StaticResponseGenerator) ResponseGenerator(c *gin.Context) (map[string]interface{}, error) {
	params := map[string]string{}
	for _, v := range c.Params {
		params[v.Key] = v.Value
	}
	target := map[string]interface{}{
		"extra":   s.Page.Extra,
		"context": c,
		"params":  params,
		"helper":  &TplHelper{},
	}
	return target, nil
}

type DynamicResponseGenerator struct {
	Page    Page
	Backend Backend
	Decoder Decoder
}

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
	target["helper"] = &TplHelper{}
	return target, nil
}

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

type TplHelper struct {
}

func (t *TplHelper) Now() string {
	return time.Now().String()
}

func NewErrorHandler(path string) (ErrorHandler, error) {
	templateFile, err := os.Open(path)
	if err != nil {
		log.Println("reading", path, ":", err.Error())
		return ErrorHandler{}, err
	}
	defer templateFile.Close()
	data, err := ioutil.ReadAll(templateFile)
	if err != nil {
		log.Println("reading", path, ":", err.Error())
		return ErrorHandler{}, err
	}
	return ErrorHandler{data}, nil
}

type ErrorHandler struct {
	Content []byte
}

func (e *ErrorHandler) HandlerFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if !c.IsAborted() {
			return
		}

		c.Writer.Write(e.Content)
	}
}

func (e *ErrorHandler) StaticHandlerFunc() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Write(e.Content)
	}
}
