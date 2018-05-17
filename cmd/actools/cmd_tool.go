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
			run := func(cmd *cobra.Command, args []string) error {
				options := []docker.ContainerOption{
					docker.WithImage(docker.Image(containers.Repo, container.Image, "latest")),
					docker.WithDefaultNetwork(),
				}
				options = append(options, container.Options...)
				container, err := docker.Container(fmt.Sprintf("tool-%s-%s", container.Image, tool), options...)
				if err != nil {
					return errors.Trace(err)
				}

				if err := container.Run(args...); err != nil {
					return errors.Trace(err)
				}

				return nil
			}

			var CmdTool = &cobra.Command{
				Use:                   tool,
				Short:                 fmt.Sprintf("Herramienta %s [%s]", tool, container.Image),
				DisableFlagParsing:    true,
				DisableFlagsInUseLine: true,
				RunE: run,
			}
			CmdRoot.AddCommand(CmdTool)

			var CmdToolDebug = &cobra.Command{
				Use:                   tool,
				Short:                 fmt.Sprintf("Herramienta %s [%s]", tool, container.Image),
				DisableFlagParsing:    true,
				DisableFlagsInUseLine: true,
				RunE: run,
			}
			CmdDebug.AddCommand(CmdToolDebug)
		}
	}
}
