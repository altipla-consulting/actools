package main

import (
	"fmt"

	"github.com/juju/errors"
	"github.com/spf13/cobra"
)

func init() {
	for tool, toolCnf := range tools {
		createToolCommand(tool, toolCnf.Container, toolCnf.Cnf)
	}
}

func createToolCommand(tool, container string, cnf *containerConfig) {
	run := func(cmd *cobra.Command, args []string) error {
		containerArgs := append([]string{tool}, args...)
		return errors.Trace(runContainer(container, cnf, containerArgs...))
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
