package main

import (
	"fmt"

	"github.com/juju/errors"
	"github.com/spf13/cobra"
)

func init() {
	for tool, container := range tools {
		createToolCommand(tool, container)
	}
}

func createToolCommand(tool, container string) {
	run := func(cmd *cobra.Command, args []string) error {
		containerArgs := append([]string{tool}, args...)
		return errors.Trace(runContainer(container, containerArgs...))
	}

	var CmdTool = &cobra.Command{
		Use:                   tool,
		Short:                 fmt.Sprintf("Herramienta %s [%s]", tool, container),
		DisableFlagParsing:    true,
		DisableFlagsInUseLine: true,
		RunE: run,
	}
	CmdRoot.AddCommand(CmdTool)

	var CmdToolDebug = &cobra.Command{
		Use:                   tool,
		Short:                 fmt.Sprintf("Herramienta %s [%s]", tool, container),
		DisableFlagParsing:    true,
		DisableFlagsInUseLine: true,
		RunE: run,
	}
	CmdDebug.AddCommand(CmdToolDebug)
}
