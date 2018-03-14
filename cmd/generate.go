package cmd

import (
	"log"
	"os"
	"time"

	"github.com/devopsfaith/api2html/generator"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
)

var (
	isos        string
	basePath    string
	ignoreRegex string

	generateCmd = &cobra.Command{
		Use:     "generate",
		Short:   "Generate the final api2html templates.",
		Long:    "Generate the final api2html templates.",
		RunE:    generatorWrapper{defaultGeneratorFactory}.Generate,
		Aliases: []string{"create", "new"},
		Example: "api2html generate -i en_US -r partial",
	}

	generateAndWatchCmd = &cobra.Command{
		Use:     "watch",
		Short:   "Generate the final api2html templates after every change in the target.",
		Long:    "Generate the final api2html templates after every change in the target.",
		RunE:    generatorWatchWrapper{generatorWrapper{defaultGeneratorFactory}}.Watch,
		Example: "api2html generate watch -i en_US -r partial",
	}
)

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.PersistentFlags().StringVarP(&basePath, "path", "p", os.Getenv("PWD"), "Base path for the generation")
	generateCmd.PersistentFlags().StringVarP(&isos, "iso", "i", "*", "(comma-separated) iso code of the site to create")
	generateCmd.PersistentFlags().StringVarP(&ignoreRegex, "reg", "r", "ignore", "regex filtering the sources to move to the output folder")

	generateCmd.AddCommand(generateAndWatchCmd)
}

type generatorFactory func(basePath string, ignoreRegex string) generator.Generator

func defaultGeneratorFactory(basePath string, ignoreRegex string) generator.Generator {
	return generator.New(basePath, ignoreRegex)
}

type generatorWrapper struct {
	gf generatorFactory
}

func (g generatorWrapper) Generate(_ *cobra.Command, _ []string) error {
	start := time.Now()

	if err := g.gf(basePath, ignoreRegex).Generate(isos); err != nil {
		log.Println("generation aborted:", err.Error())
		return err
	}

	log.Println("site generated! time:", time.Since(start))
	return nil
}

type generatorWatchWrapper struct {
	generatorWrapper
}

func (g generatorWatchWrapper) Watch(c *cobra.Command, p []string) error {
	if err := g.Generate(c, p); err != nil {
		return err
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	err = watcher.Add(basePath)
	if err != nil {
		return err
	}

	for {
		select {
		case event := <-watcher.Events:
			log.Println("event:", event)
			if err := g.Generate(c, p); err != nil {
				return err
			}
		case err := <-watcher.Errors:
			return err
		}
	}
}
