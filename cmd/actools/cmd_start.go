package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/altipla-consulting/collections"
	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/altipla-consulting/actools/pkg/config"
	"github.com/altipla-consulting/actools/pkg/containers"
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
			docker.WithImage(docker.Image(containers.Repo, containerDesc.Image, "latest")),
			docker.WithDefaultNetwork(),
			docker.WithPersistence(),
			docker.WithNetworkAlias(tool),
		}
		options = append(options, containerDesc.Options...)

		for _, port := range config.Settings.Tools[tool].Ports {
			parts := strings.Split(port, ":")
			if len(parts) != 2 {
				return errors.NotValidf("ports of tool %s", tool)
			}

			source, err := strconv.ParseInt(parts[0], 10, 64)
			if err != nil {
				return errors.NewNotValid(err, fmt.Sprintf("invalid port number: %s", parts[0]))
			}

			inside, err := strconv.ParseInt(parts[1], 10, 64)
			if err != nil {
				return errors.NewNotValid(err, fmt.Sprintf("invalid port number: %s", parts[1]))
			}

			options = append(options, docker.WithPort(source, inside))
		}

		for _, volume := range config.Settings.Tools[tool].Volumes {
			parts := strings.Split(volume, ":")
			if len(parts) != 2 {
				return errors.NotValidf("volumes of tool %s", tool)
			}

			options = append(options, docker.WithVolume(parts[0], parts[1]))
		}

		container, err := docker.Container(tool, options...)
		if err != nil {
			return errors.Trace(err)
		}

		log.WithField("service", tool).Info("Start service")
		if err := container.Start(config.Settings.Tools[tool].Args...); err != nil {
			return errors.Trace(err)
		}
	}

	watcher := docker.NewWatcher()
	for _, service := range services {
		containerDesc, err := containers.FindImage(fmt.Sprintf("dev-%s", config.Settings.Services[service].Type))
		if err != nil {
			return errors.Trace(err)
		}

		options := []docker.ContainerOption{
			docker.WithImage(docker.Image(containers.Repo, containerDesc.Image, "latest")),
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
			parts := strings.Split(port, ":")
			if len(parts) != 2 {
				return errors.NotValidf("ports of service %s", service)
			}

			source, err := strconv.ParseInt(parts[0], 10, 64)
			if err != nil {
				return errors.NewNotValid(err, fmt.Sprintf("invalid port number: %s", parts[0]))
			}

			inside, err := strconv.ParseInt(parts[1], 10, 64)
			if err != nil {
				return errors.NewNotValid(err, fmt.Sprintf("invalid port number: %s", parts[1]))
			}

			options = append(options, docker.WithPort(source, inside))
		}

		for _, volume := range config.Settings.Services[service].Volumes {
			parts := strings.Split(volume, ":")
			if len(parts) != 2 {
				return errors.NotValidf("volumes of service %s", service)
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

	watcher.Wait()

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

		return nil, nil, errors.NotFoundf("service %s", arg)
	}

	services = collections.UniqueStrings(services)
	tools = collections.UniqueStrings(tools)

	return services, tools, nil
}
