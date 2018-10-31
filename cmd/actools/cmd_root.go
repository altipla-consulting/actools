package main

import (
	"os"

	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/altipla-consulting/actools/pkg/update"
)

var debugApp bool

func init() {
	CmdRoot.PersistentFlags().BoolVarP(&debugApp, "debug", "d", false, "Activa el logging de depuración")
}

var CmdRoot = &cobra.Command{
	Use:          "actools",
	Short:        "Actools actua de helper de las herramientas de desarrollo de Altipla Consulting",
	SilenceUsage: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if debugApp {
			log.SetLevel(log.DebugLevel)
			log.Debug("DEBUG log level activated")
		}

		if err := update.Check(); err != nil {
			return errors.Trace(err)
		}

		return nil
	},
}
