package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	CmdRoot.AddCommand(CmdUpdate)
}

var CmdUpdate = &cobra.Command{
	Use:   "update",
	Short: "Imprime el comando de actualización de la herramienta.",
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Info()
		log.Info("Run the following command to install the latest version:")
		log.Info()
		log.Info("\tcurl -L https://git.io/Jvihi | bash")
		log.Info()

		return nil
	},
}
