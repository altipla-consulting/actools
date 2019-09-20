package main

import (
	"github.com/spf13/cobra"
	"libs.altipla.consulting/errors"
)

func init() {
	CmdRoot.AddCommand(CmdRestart)
}

var CmdRestart = &cobra.Command{
	Use:   "restart",
	Short: "Reinicia un servicio de desarrollo.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.Errorf("arguments required")
		}

		if err := stopCommand(args); err != nil {
			return errors.Trace(err)
		}
		if err := startCommand(args); err != nil {
			return errors.Trace(err)
		}

		return nil
	},
}
