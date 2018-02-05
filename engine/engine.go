package engine

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func New(cfgPath string, devel bool) (*gin.Engine, error) {
	cfg, err := ParseConfigFromFile(cfgPath)
	if err != nil {
		return nil, err
	}

	templateStore := NewTemplateStore()
	if !devel {
		gin.SetMode(gin.ReleaseMode)
	}
	e := gin.Default()
	e.RedirectTrailingSlash = true
	e.RedirectFixedPath = true

	setStatics(e, cfg)

	pf := NewMustachePageFactory(e, templateStore)
	pf.Build(cfg)

	if h, err := NewStaticHandler("./static/404"); err == nil {
		e.NoRoute(h.HandlerFunc())
	}

	if h, err := NewStaticHandler("./static/405"); err == nil {
		e.NoMethod(h.HandlerFunc())
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

func setStatics(e *gin.Engine, cfg Config) {
	if h, err := NewErrorHandler("./static/500"); err == nil {
		e.Use(h.HandlerFunc())
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
}
