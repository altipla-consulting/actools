package main

import (
	"fmt"

	"github.com/juju/errors"
	"github.com/spf13/cobra"
)

func init() {
	CmdRoot.AddCommand(CmdRun)
	for _, container := range containers {
		createRunCommand(container)
	}
}

var CmdRun = &cobra.Command{
	Use:   "run",
	Short: "Ejecuta un comando manualmente dentro de un contenedor de herramienta",
}

func createRunCommand(container string) {
	var CmdContainer = &cobra.Command{
		Use:                   container,
		Short:                 fmt.Sprintf("Contenedor %s", container),
		DisableFlagParsing:    true,
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return errors.Trace(runContainer(container, nil, args...))
		},
	}
	CmdRun.AddCommand(CmdContainer)
}
