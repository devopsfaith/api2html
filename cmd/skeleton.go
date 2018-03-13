package cmd

import (
	"github.com/devopsfaith/api2html/skeleton"
	"github.com/spf13/cobra"
)

var (
	outputPath string

	skelCmd = &cobra.Command{
		Use:     "skeleton",
		Short:   "Run the api2html server.",
		Long:    "Run the api2html server.",
		RunE:    skelWrapper{defaultSkelFactory}.Create,
		Aliases: []string{"skel"},
		Example: "api2html skeleton -o skel_example",
	}
)

func init() {
	rootCmd.AddCommand(skelCmd)

	skelCmd.PersistentFlags().StringVarP(&outputPath, "outputPath", "o", "skel_example", "Output path for the skel generation")
}

type skelFactory func(outputPath string) skeleton.Skel

func defaultSkelFactory(outputPath string) skeleton.Skel {
	return skeleton.New(outputPath)
}

type skelWrapper struct {
	sk skelFactory
}

func (sw skelWrapper) Create(_ *cobra.Command, _ []string) error {
	return sw.sk(outputPath).Create()
}
