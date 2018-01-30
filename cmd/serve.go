package cmd

import (
	"log"
	"time"

	"github.com/devopsfaith/api2html/engine"
	"github.com/spf13/cobra"
)

var (
	cfgFile string
	devel   bool

	serveCmd = &cobra.Command{
		Use:     "serve",
		Short:   "Run the api2html server.",
		Long:    "Run the api2html server.",
		Run:     serve,
		Aliases: []string{"run", "server", "start"},
		Example: "api2html serve -d -c config.json",
	}
)

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "config.json", "Path to the configuration filename")
	serveCmd.PersistentFlags().BoolVarP(&devel, "devel", "d", false, "Enable the devel")
}

func serve(cmd *cobra.Command, args []string) {
	engine, err := engine.New(cfgFile, devel)
	if err != nil {
		log.Println("engine creation aborted:", err.Error())
		return
	}

	time.Sleep(time.Second)

	engine.Run(":8080")
}
