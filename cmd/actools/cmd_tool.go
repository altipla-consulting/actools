package main

import (
	"fmt"

	"github.com/juju/errors"
	"github.com/spf13/cobra"

	"github.com/altipla-consulting/actools/pkg/containers"
	"github.com/altipla-consulting/actools/pkg/docker"
)

func init() {
	for _, container := range containers.List() {
		for _, tool := range container.Tools {
			var CmdTool = &cobra.Command{
				Use:                   tool,
				Short:                 fmt.Sprintf("Herramienta %s [%s]", tool, container.Image),
				DisableFlagParsing:    true,
				DisableFlagsInUseLine: true,
				RunE: createToolEntrypoint(container, tool),
			}
			CmdRoot.AddCommand(CmdTool)

			var CmdToolDebug = &cobra.Command{
				Use:                   tool,
				Short:                 fmt.Sprintf("Herramienta %s [%s]", tool, container.Image),
				DisableFlagParsing:    true,
				DisableFlagsInUseLine: true,
				RunE: createToolEntrypoint(container, tool),
			}
			CmdDebug.AddCommand(CmdToolDebug)
		}
	}
}

func createToolEntrypoint(containerDesc containers.Container, tool string) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		options := []docker.ContainerOption{
			docker.WithImage(docker.Image(containers.Repo, containerDesc.Image, "latest")),
			docker.WithDefaultNetwork(),
		}
		options = append(options, containerDesc.Options...)

		container, err := docker.Container(fmt.Sprintf("tool-%s-%s", containerDesc.Image, tool), options...)
		if err != nil {
			return errors.Trace(err)
		}

		args = append([]string{tool}, args...)

		if err := container.Run(args...); err != nil {
			return errors.Trace(err)
		}

		return nil
	}
}
