package cmd

import (
	"github.com/devopsfaith/api2html/skeleton"
	"github.com/spf13/cobra"
)

var (
	outputPath string

	skelCmd = &cobra.Command{
		Use:     "skel",
		Short:   "Executes commands to manage/create skeletons",
		Aliases: []string{"skeleton"},
		Example: "api2html skel",
	}

	createCmd = &cobra.Command{
		Use:     "create",
		Short:   "Creates a skeleton structure.",
		Example: "api2html skel create -o outputPath <skelname>",
	}

	blogCmd = &cobra.Command{
		Use:     "blog",
		Short:   "Creates the blog skeleton example.",
		RunE:    skelWrapper{defaultBlogSkelFactory}.Create,
		Example: "api2html skel create -o outputPath blog",
	}
)

func init() {
	rootCmd.AddCommand(skelCmd)
	skelCmd.AddCommand(createCmd)
	createCmd.AddCommand(blogCmd)

	createCmd.PersistentFlags().StringVarP(&outputPath, "outputPath", "o", "example", "Output path for the example generation skel")

}

type skelFactory func(outputPath string) skeleton.Skel

func defaultBlogSkelFactory(outputPath string) skeleton.Skel {
	return skeleton.NewBlog(outputPath)
}

type skelWrapper struct {
	sk skelFactory
}

func (sw skelWrapper) Create(_ *cobra.Command, _ []string) error {
	return sw.sk(outputPath).Create()
}
