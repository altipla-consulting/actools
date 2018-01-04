package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var debugApp bool

func init() {
	CmdRoot.PersistentFlags().BoolVarP(&debugApp, "debug", "d", false, "Activa el logging de depuración")
}

var CmdRoot = &cobra.Command{
	Use:   "actools",
	Short: "Actools actua de helper de las herramientas de desarrollo de Altipla Consulting",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if debugApp {
			log.SetLevel(log.DebugLevel)
			log.Debug("DEBUG log level activated")
		}
	},
}
