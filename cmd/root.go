package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "api2html",
	Short: "Generate HTML on the fly from your API.",
}

// Execute executes the root command
func Execute() error {
	return rootCmd.Execute()
}
