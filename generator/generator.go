package generator

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"
)

type Generator interface {
	Generate(isos string) error
}

func New(basePath, ignoreRegex string) Generator {
	return &BasicGenerator{
		SourceFolder:     fmt.Sprintf("%s/sources", basePath),
		ConfigFolder:     fmt.Sprintf("%s/config", basePath),
		I18NFolder:       fmt.Sprintf("%s/i18n", basePath),
		OutputFolder:     fmt.Sprintf("%s/output", basePath),
		IgnorePattern:    ignoreRegex,
		ScannerFactory:   NewScanner,
		CollectorFactory: NewCollector,
		RendererFactory:  NewRenderer,
	}
}

type BasicGenerator struct {
	SourceFolder     string
	I18NFolder       string
	ConfigFolder     string
	OutputFolder     string
	IgnorePattern    string
	ScannerFactory   func([]string) Scanner
	CollectorFactory func(string, string) Collector
	RendererFactory  func(string, *regexp.Regexp) Renderer
}

func (g *BasicGenerator) Generate(isos string) error {
	collector := g.CollectorFactory(g.ConfigFolder, g.I18NFolder)
	renderer := g.RendererFactory(g.OutputFolder, regexp.MustCompile(g.IgnorePattern))

	if isos == "*" {
		isos = strings.Join(collector.AvailableISOs(), ",")
	}

	for _, iso := range strings.Split(isos, ",") {
		start := time.Now()
		log.Printf("[%s] generating the site", iso)

		data, err := collector.Collect(iso)
		if err != nil {
			log.Println(err.Error())
			return err
		}

		log.Printf("[%s] all translations and configurations collected", iso)
		log.Printf("[%s] %v", iso, data)

		scanner := g.ScannerFactory([]string{
			fmt.Sprintf("%s/global", g.SourceFolder),
			fmt.Sprintf("%s/%s", g.SourceFolder, iso),
		})

		if err := renderer.Render(iso, data, scanner); err != nil {
			log.Printf("[%s] error: %s", iso, err.Error())
			return err
		}

		log.Println("****************************************")
		log.Printf("[%s] site generated! time: %s", iso, time.Since(start).String())
		log.Println("****************************************")
	}
	return nil
}
