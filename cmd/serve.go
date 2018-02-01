package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/devopsfaith/api2html/engine"
	"github.com/spf13/cobra"
)

var (
	cfgFile string
	devel   bool
	port    int

	serveCmd = &cobra.Command{
		Use:     "serve",
		Short:   "Run the api2html server.",
		Long:    "Run the api2html server.",
		RunE:    serveWrapper{defaultEngineFactory}.Serve,
		Aliases: []string{"run", "server", "start"},
		Example: "api2html serve -d -c config.json",
	}

	errNilEngine = fmt.Errorf("serve cmd aborted: nil engine")
)

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "config.json", "Path to the configuration filename")
	serveCmd.PersistentFlags().BoolVarP(&devel, "devel", "d", false, "Enable the devel")
	serveCmd.PersistentFlags().IntVarP(&port, "port", "p", 8080, "Listen port")
}

type engineWrapper interface {
	Run(...string) error
}

type engineFactory func(cfgPath string, devel bool) (engineWrapper, error)

func defaultEngineFactory(cfgPath string, devel bool) (engineWrapper, error) {
	return engine.New(cfgPath, devel)
}

type serveWrapper struct {
	eF engineFactory
}

func (s serveWrapper) Serve(_ *cobra.Command, _ []string) error {
	eW, err := s.eF(cfgFile, devel)
	if err != nil {
		log.Println("engine creation aborted:", err.Error())
		return err
	}
	if eW == nil {
		log.Println("engine creation aborted:", errNilEngine.Error())
		return errNilEngine
	}

	time.Sleep(time.Second)

	return eW.Run(fmt.Sprintf(":%d", port))
}
