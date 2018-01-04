package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var CmdDebug = &cobra.Command{
	Use:   "debug",
	Short: "Activa el modo depuración de las herramientas",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		log.SetLevel(log.DebugLevel)
		log.Debug("DEBUG log level activated")
	},
}

func init() {
	CmdRoot.AddCommand(CmdDebug)
}
