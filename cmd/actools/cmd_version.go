package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/altipla-consulting/actools/pkg/config"
)

func init() {
	CmdRoot.AddCommand(CmdVersion)
}

var CmdVersion = &cobra.Command{
	Use:   "version",
	Short: "Imprime la versión de la herramienta.",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(config.Version)

		return nil
	},
}
