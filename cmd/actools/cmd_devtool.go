package main

import (
	"github.com/spf13/cobra"
	"libs.altipla.consulting/errors"

	"github.com/altipla-consulting/actools/pkg/devtool"
)

func init() {
	CmdRoot.AddCommand(CmdDevTool)
}

var CmdDevTool = &cobra.Command{
	Use:   "devtool",
	Short: "Inicia el servicio de desarrollo de aplicaciones en Go.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.Trace(devtool.Run(args))
	},
}
