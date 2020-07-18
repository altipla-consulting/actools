package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"libs.altipla.consulting/errors"

	"github.com/altipla-consulting/actools/pkg/config"
	"github.com/altipla-consulting/actools/pkg/containers"
	"github.com/altipla-consulting/actools/pkg/docker"
)

func init() {
	for _, container := range containers.List() {
		for _, tool := range container.Tools {
			var CmdToolDirect = &cobra.Command{
				Use:                   tool,
				Short:                 fmt.Sprintf("Herramienta %s [%s]", tool, container.Image),
				DisableFlagParsing:    true,
				DisableFlagsInUseLine: true,
				RunE:                  createToolEntrypoint(container, tool, ""),
			}
			CmdRoot.AddCommand(CmdToolDirect)

			var CmdToolDebug = &cobra.Command{
				Use:                   tool,
				Short:                 fmt.Sprintf("Herramienta %s [%s]", tool, container.Image),
				DisableFlagParsing:    true,
				DisableFlagsInUseLine: true,
				RunE:                  createToolEntrypoint(container, tool, ""),
			}
			CmdDebug.AddCommand(CmdToolDebug)
		}
	}
}

func createToolEntrypoint(containerDesc containers.Container, tool, workdir string) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		options := []docker.ContainerOption{
			docker.WithImage(docker.Image(containers.Repo, containerDesc.Image)),
			docker.WithDefaultNetwork(),
			docker.WithEnv("PROJECT", config.Settings.Project),

			// Useful mostly for Jenkins.
			docker.WithEnv("BUILD_NUMBER", os.Getenv("BUILD_NUMBER")),
		}
		options = append(options, containerDesc.Options...)
		if workdir != "" {
			options = append(options, docker.WithWorkdir(fmt.Sprintf("/workspace/%s", workdir)))
		}

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
