package cmd

import (
	"log"
	"os"
	"time"

	"github.com/devopsfaith/api2html/engine"
	"github.com/devopsfaith/api2html/generator"
	"github.com/spf13/cobra"
)

var (
	cfgFile string
	devel   bool

	isos        string
	basePath    string
	ignoreRegex string

	rootCmd = &cobra.Command{
		Use:   "api2html",
		Short: "Template Render As A Service",
	}

	generateCmd = &cobra.Command{
		Use:   "generate",
		Short: "Generate the final api2html templates.",
		Long:  "Generate the final api2html templates.",
		Run: func(cmd *cobra.Command, args []string) {
			start := time.Now()

			if err := generator.New(basePath, ignoreRegex).Generate(isos); err != nil {
				log.Println("generation aborted:", err.Error())
				return
			}

			log.Println("site generated! time:", time.Since(start))
		},
		Aliases: []string{"create", "new"},
		Example: "api2html generate -d -c config.json",
	}

	serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "Run the api2html server.",
		Long:  "Run the api2html server.",
		Run: func(cmd *cobra.Command, args []string) {
			engine, err := engine.New(cfgFile, devel)
			if err != nil {
				log.Println("engine creation aborted:", err.Error())
				return
			}

			time.Sleep(time.Second)

			engine.Run(":8080")
		},
		Aliases: []string{"run", "server", "start"},
		Example: "api2html serve -d -c config.json",
	}
)

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.PersistentFlags().StringVarP(&basePath, "path", "p", os.Getenv("PWD"), "Base path for the generation")
	generateCmd.PersistentFlags().StringVarP(&isos, "iso", "i", "*", "(comma-separated) iso code of the site to create")
	generateCmd.PersistentFlags().StringVarP(&ignoreRegex, "reg", "r", "ignore", "regex filtering the sources to move to the output folder")

	rootCmd.AddCommand(serveCmd)
	serveCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "config.json", "Path to the configuration filename")
	serveCmd.PersistentFlags().BoolVarP(&devel, "devel", "d", false, "Enable the devel")
}

func Execute() error {
	return rootCmd.Execute()
}
