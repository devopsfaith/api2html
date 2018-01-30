package cmd

import (
	"log"
	"os"
	"time"

	"github.com/devopsfaith/api2html/generator"
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
		Run:     generate,
		Aliases: []string{"create", "new"},
		Example: "api2html generate -i en_US -r partial",
	}
)

func init() {
	rootCmd.AddCommand(generateCmd)

	generateCmd.PersistentFlags().StringVarP(&basePath, "path", "p", os.Getenv("PWD"), "Base path for the generation")
	generateCmd.PersistentFlags().StringVarP(&isos, "iso", "i", "*", "(comma-separated) iso code of the site to create")
	generateCmd.PersistentFlags().StringVarP(&ignoreRegex, "reg", "r", "ignore", "regex filtering the sources to move to the output folder")
}

func generate(cmd *cobra.Command, args []string) {
	start := time.Now()

	if err := generator.New(basePath, ignoreRegex).Generate(isos); err != nil {
		log.Println("generation aborted:", err.Error())
		return
	}

	log.Println("site generated! time:", time.Since(start))
}
