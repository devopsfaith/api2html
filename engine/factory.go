package engine

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	newrelic "github.com/newrelic/go-agent"
	nrgin "github.com/newrelic/go-agent/_integrations/nrgin/v1"
)

// DefaultFactory is an Factory ready to be used
var DefaultFactory = Factory{
	TemplateStoreFactory: NewTemplateStore,
	Parser:               ParseConfigFromFile,
	MustachePageFactory:  NewMustachePageFactory,
	StaticHandlerFactory: NewStaticHandler,
	ErrorHandlerFactory:  NewErrorHandler,
}

// Factory is a struct able to build api2html engines
type Factory struct {
	TemplateStoreFactory func() *TemplateStore
	Parser               func(string) (Config, error)
	MustachePageFactory  func(*gin.Engine, *TemplateStore) MustachePageFactory
	StaticHandlerFactory func(string) (StaticHandler, error)
	ErrorHandlerFactory  func(string, int) (ErrorHandler, error)
}

// New creates a gin engine with the received config and the injected factories
func (ef Factory) New(cfgPath string, devel bool) (*gin.Engine, error) {
	cfg, err := ef.Parser(cfgPath)
	if err != nil {
		return nil, err
	}

	if cfg.NewRelic != nil && cfg.NewRelic.License != "" {
		nrCfg := newrelic.NewConfig(cfg.NewRelic.AppName, cfg.NewRelic.License)
		if devel {
			nrCfg.Logger = newrelic.NewDebugLogger(os.Stdout)
		}
		nrapp, err := newrelic.NewApplication(nrCfg)
		if err != nil {
			return nil, err
		}
		newrelicApp = &nrapp
	}

	templateStore := ef.TemplateStoreFactory()
	e := ef.newGinEngine(cfg, devel)
	pf := ef.MustachePageFactory(e, templateStore)
	pf.Build(cfg)

	if h, err := ef.StaticHandlerFactory("./static/404"); err == nil {
		e.NoRoute(h.HandlerFunc())
	} else {
		log.Println("using the default 404 template")
		e.NoRoute(Default404StaticHandler.HandlerFunc())
	}

	if devel {
		e.PUT("/template/:templateName", func(c *gin.Context) {
			file, err := c.FormFile("file")
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}

			f, err := file.Open()
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}

			defer f.Close()

			tmp, err := NewMustacheRenderer(f)
			if err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}

			templateName := c.Param("templateName")
			if err := templateStore.Set(templateName, tmp); err != nil {
				c.AbortWithError(http.StatusInternalServerError, err)
				return
			}

			c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded and stored as [%s]!", templateName, file.Filename))
		})
	}
	return e, nil
}

func (ef Factory) newGinEngine(cfg Config, devel bool) *gin.Engine {
	if !devel {
		gin.SetMode(gin.ReleaseMode)
	}
	e := gin.Default()
	e.RedirectTrailingSlash = true
	e.RedirectFixedPath = true

	if newrelicApp != nil {
		e.Use(nrgin.Middleware(*newrelicApp))
	}
	ef.setStatics(e, cfg)

	return e
}

func (ef Factory) setStatics(e *gin.Engine, cfg Config) {
	if cfg.PublicFolder != nil {
		e.Use(static.Serve(cfg.PublicFolder.Prefix, static.LocalFile(cfg.PublicFolder.Path, false)))
	}

	if cfg.Robots {
		log.Println("registering the robots file")
		e.StaticFile("/robots.txt", "./static/robots.txt")
	}

	if cfg.Sitemap {
		log.Println("registering the sitemap file")
		e.StaticFile("/sitemap.xml", "./static/sitemap.xml")
	}

	for _, fileName := range cfg.StaticTXTContent {
		log.Println("registering the static", fileName)
		e.StaticFile(fmt.Sprintf("/%s", fileName), fmt.Sprintf("./static/%s", fileName))
	}

	if h, err := ef.ErrorHandlerFactory("./static/500", http.StatusInternalServerError); err == nil {
		e.Use(h.HandlerFunc())
	} else {
		log.Println("using the default 500 template")
		e.Use(Default500StaticHandler.HandlerFunc())
	}

}
