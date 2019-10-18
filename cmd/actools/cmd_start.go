package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"

	"github.com/altipla-consulting/collections"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
	"libs.altipla.consulting/errors"

	"github.com/altipla-consulting/actools/pkg/config"
	"github.com/altipla-consulting/actools/pkg/containers"
	"github.com/altipla-consulting/actools/pkg/devtool"
	"github.com/altipla-consulting/actools/pkg/docker"
)

func init() {
	CmdRoot.AddCommand(CmdStart)
}

var CmdStart = &cobra.Command{
	Use:   "start",
	Short: "Enciende un servicio de desarrollo",
	RunE: func(cmd *cobra.Command, args []string) error {
		return errors.Trace(startCommand(args))
	},
}

func startCommand(args []string) error {
	services, tools, err := resolveDeps(args)
	if err != nil {
		return errors.Trace(err)
	}

	for _, tool := range tools {
		containerDesc, err := containers.FindImage(config.Settings.Tools[tool].Container)
		if err != nil {
			return errors.Trace(err)
		}

		options := []docker.ContainerOption{
			docker.WithImage(docker.Image(containers.Repo, containerDesc.Image)),
			docker.WithDefaultNetwork(),
			docker.WithPersistence(),
			docker.WithNetworkAlias(tool),
		}
		options = append(options, containerDesc.Options...)

		for _, port := range config.Settings.Tools[tool].Ports {
			options = append(options, docker.WithPorts(port))
		}

		for _, volume := range config.Settings.Tools[tool].Volumes {
			parts := strings.Split(volume, ":")
			if len(parts) != 2 {
				return errors.Errorf("invalid volumes of tool: %s", tool)
			}

			options = append(options, docker.WithVolume(parts[0], parts[1]))
		}

		container, err := docker.Container(tool, options...)
		if err != nil {
			return errors.Trace(err)
		}

		log.WithField("tool", tool).Info("Start tool")
		if err := container.Start(config.Settings.Tools[tool].Args...); err != nil {
			return errors.Trace(err)
		}
	}

	if len(services) > 0 {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		g, ctx := errgroup.WithContext(ctx)

		watcher := docker.NewWatcher()
		g.Go(func() error {
			<-ctx.Done()
			return errors.Trace(watcher.StopAll())
		})

		var goapps []string
		for _, service := range services {
			if config.Settings.Services[service].Type == "go" {
				goapps = append(goapps, service)
				continue
			}

			containerDesc, err := containers.FindImage(fmt.Sprintf("dev-%s", config.Settings.Services[service].Type))
			if err != nil {
				return errors.Trace(err)
			}

			options := []docker.ContainerOption{
				docker.WithImage(docker.Image(containers.Repo, containerDesc.Image)),
				docker.WithDefaultNetwork(),
				docker.WithPersistence(),
				docker.WithWorkdir(fmt.Sprintf("/workspace/%s", config.Settings.Services[service].Workdir)),
				docker.WithNetworkAlias(service),
				docker.WithEnv("PROJECT", config.Settings.Project),
				docker.WithEnv("WORKDIR", config.Settings.Services[service].Workdir),
				docker.WithEnv("SERVICE", service),
			}
			options = append(options, containerDesc.Options...)

			for _, port := range config.Settings.Services[service].Ports {
				options = append(options, docker.WithPorts(port))
			}

			for _, volume := range config.Settings.Services[service].Volumes {
				parts := strings.Split(volume, ":")
				if len(parts) != 2 {
					return errors.Errorf("invalid volumes of service: %s", service)
				}

				options = append(options, docker.WithVolume(parts[0], parts[1]))
			}

			for k, v := range config.Settings.Services[service].Env {
				options = append(options, docker.WithEnv(k, v))
			}

			container, err := docker.Container(service, options...)
			if err != nil {
				return errors.Trace(err)
			}

			watcher.Run(service, container)
		}
		if len(goapps) > 0 {
			if err := devtool.Run(ctx, g, watcher, goapps); err != nil {
				return errors.Trace(err)
			}
		}

		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt)
		for range c {
			// Newline to always jump the next log we emit.
			fmt.Println()
			cancel()
			break
		}
		return errors.Trace(g.Wait())
	}

	return nil
}

func resolveDeps(args []string) ([]string, []string, error) {
	if len(args) == 0 {
		return nil, nil, nil
	}

	services := []string{}
	tools := []string{}
	for _, arg := range args {
		if config.Settings.IsService(arg) {
			services = append(services, arg)

			subservices, subtools, err := resolveDeps(config.Settings.Services[arg].Deps)
			if err != nil {
				return nil, nil, errors.Trace(err)
			}
			services = append(services, subservices...)
			tools = append(tools, subtools...)

			continue
		}

		if config.Settings.IsTool(arg) {
			tools = append(tools, arg)

			subservices, subtools, err := resolveDeps(config.Settings.Tools[arg].Deps)
			if err != nil {
				return nil, nil, errors.Trace(err)
			}
			services = append(services, subservices...)
			tools = append(tools, subtools...)

			continue
		}

		return nil, nil, errors.Errorf("service not found: %s", arg)
	}

	services = collections.UniqueStrings(services)
	tools = collections.UniqueStrings(tools)

	return services, tools, nil
}
