package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"libs.altipla.consulting/errors"

	"github.com/altipla-consulting/actools/pkg/containers"
	"github.com/altipla-consulting/actools/pkg/docker"
)

func init() {
	CmdRoot.AddCommand(CmdRun)
	for _, container := range containers.List() {
		var CmdContainer = &cobra.Command{
			Use:                   container.Image,
			Short:                 fmt.Sprintf("Contenedor %s", container.Image),
			DisableFlagParsing:    true,
			DisableFlagsInUseLine: true,
			RunE:                  createRunEntrypoint(container),
		}
		CmdRun.AddCommand(CmdContainer)
	}
}

var CmdRun = &cobra.Command{
	Use:   "run",
	Short: "Ejecuta un comando manualmente dentro de un contenedor de herramienta concreto.",
}

func createRunEntrypoint(containerDesc containers.Container) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		options := []docker.ContainerOption{
			docker.WithImage(docker.Image(containers.Repo, containerDesc.Image)),
			docker.WithDefaultNetwork(),
		}
		options = append(options, containerDesc.Options...)

		container, err := docker.Container(fmt.Sprintf("run-%s", containerDesc.Image), options...)
		if err != nil {
			return errors.Trace(err)
		}

		if err := container.Run(args...); err != nil {
			return errors.Trace(err)
		}

		return nil
	}
}
