package main

import (
	"fmt"

	"github.com/juju/errors"
	"github.com/spf13/cobra"

	"github.com/altipla-consulting/actools/pkg/config"
	"github.com/altipla-consulting/actools/pkg/containers"
	"github.com/altipla-consulting/actools/pkg/docker"
)

func init() {
	var CmdApp = &cobra.Command{
		Use:   "app",
		Short: "Ejecuta una herramienta dentro de la carpeta de una aplicación",
	}
	CmdRoot.AddCommand(CmdApp)

	CmdsApp := map[string]*cobra.Command{}
	for name := range config.Settings.Services {
		var CmdAppService = &cobra.Command{
			Use:   name,
			Short: fmt.Sprintf("Servicio %s", name),
		}
		CmdApp.AddCommand(CmdAppService)

		CmdsApp[name] = CmdAppService
	}

	for _, container := range containers.List() {
		for _, tool := range container.Tools {
			var CmdToolDirect = &cobra.Command{
				Use:                   tool,
				Short:                 fmt.Sprintf("Herramienta %s [%s]", tool, container.Image),
				DisableFlagParsing:    true,
				DisableFlagsInUseLine: true,
				RunE: createToolEntrypoint(container, tool, ""),
			}
			CmdRoot.AddCommand(CmdToolDirect)

			var CmdToolDebug = &cobra.Command{
				Use:                   tool,
				Short:                 fmt.Sprintf("Herramienta %s [%s]", tool, container.Image),
				DisableFlagParsing:    true,
				DisableFlagsInUseLine: true,
				RunE: createToolEntrypoint(container, tool, ""),
			}
			CmdDebug.AddCommand(CmdToolDebug)

			for name, service := range config.Settings.Services {
				var CmdAppServiceTool = &cobra.Command{
					Use:                   tool,
					Short:                 fmt.Sprintf("Herramienta %s [%s]", tool, container.Image),
					DisableFlagParsing:    true,
					DisableFlagsInUseLine: true,
					RunE: createToolEntrypoint(container, tool, service.Workdir),
				}
				CmdsApp[name].AddCommand(CmdAppServiceTool)
			}
		}
	}
}

func createToolEntrypoint(containerDesc containers.Container, tool, workdir string) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		options := []docker.ContainerOption{
			docker.WithImage(docker.Image(containers.Repo, containerDesc.Image)),
			docker.WithDefaultNetwork(),
			docker.WithEnv("PROJECT", config.Settings.Project),
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
