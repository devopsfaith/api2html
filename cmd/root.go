package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "api2html",
	Short: "Template Render As A Service",
}

// Execute executes the root command
func Execute() error {
	return rootCmd.Execute()
}
