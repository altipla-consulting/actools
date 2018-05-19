package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"

	"github.com/altipla-consulting/collections"
	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/altipla-consulting/actools/pkg/docker"
)

func init() {
	CmdRoot.AddCommand(CmdStart)
}

var CmdStart = &cobra.Command{
	Use:   "start",
	Short: "Enciende un servicio de desarrollo",
	RunE: func(cmd *cobra.Command, args []string) error {
		cnf, err := ReadConfig()
		if err != nil {
			return errors.Trace(err)
		}
		if cnf == nil {
			return errors.NotFoundf("actools.yml")
		}

		services, tools, err := resolveDeps(cnf, args)
		if err != nil {
		  return errors.Trace(err)
		}

		for _, tool := range tools {
			containerDesc := containers.FindImage(cnf.Tools[tool].Container)

			options := []docker.ContainerOption{
				docker.WithImage(docker.Image(containers.Repo, containerDesc.Image, "latest")),
				docker.WithDefaultNetwork(),
			}
			options = append(options, containerDesc.Options...)

			for _, port := range cnf.Tools[tool].Ports {
				parts := strings.Split(port, ":")
				if len(parts) != 2 {
					return errors.NotValidf("ports of tool %s", tool)
				}

				source, err := strconv.ParseInt(parts[0], 10, 64)
				if err != nil {
				  return errors.NewNotValid(err)
				}

				inside, err := strconv.ParseInt(parts[1], 10, 64)
				if err != nil {
				  return errors.NewNotValid(err)
				}

				options = append(options, docker.WithPort(source, inside))
			}

			for _, volumes := range cnf.Tools[tool].Volumes {
				parts := strings.Split(port, ":")
				if len(parts) != 2 {
					return errors.NotValidf("volumes of tool %s", tool)
				}

				options = append(options, docker.WithVolume(parts[0], parts[1]))
			}

			container, err := docker.Container(fmt.Sprintf("tool-%s", tool), options...)
			if err != nil {
				return errors.Trace(err)
			}

			log.WithField("service", tool).Info("Start service")
			if err := container.Start(cnf.Tools[tool].Args...); err != nil {
			  return errors.Trace(err)
			}
		}

		watcher := docker.NewWatcher()
		for _, service := range services {
			containerDesc := containers.FindImage(fmt.Sprintf("dev-%s", cnf.Services[service].Container))

			options := []docker.ContainerOption{
				docker.WithImage(docker.Image(containers.Repo, containerDesc.Image, "latest")),
				docker.WithDefaultNetwork(),
				docker.WithWorkdir(fmt.Sprintf("/workspace/%s", cnf.Services[service].Workdir))
			}
			options = append(options, containerDesc.Options...)

			for _, port := range cnf.Services[service].Ports {
				parts := strings.Split(port, ":")
				if len(parts) != 2 {
					return errors.NotValidf("ports of service %s", service)
				}

				source, err := strconv.ParseInt(parts[0], 10, 64)
				if err != nil {
				  return errors.NewNotValid(err)
				}

				inside, err := strconv.ParseInt(parts[1], 10, 64)
				if err != nil {
				  return errors.NewNotValid(err)
				}

				options = append(options, docker.WithPort(source, inside))
			}

			for _, volumes := range cnf.Services[service].Volumes {
				parts := strings.Split(port, ":")
				if len(parts) != 2 {
					return errors.NotValidf("volumes of service %s", service)
				}

				options = append(options, docker.WithVolume(parts[0], parts[1]))
			}

			for k, v := range env {
				options = append(options, docker.WithEnv(k, v))
			}

			container, err := docker.Container(fmt.Sprintf("service-%s", service), options...)
			if err != nil {
				return errors.Trace(err)
			}

			if err := watcher.Run(service, container); err != nil {
			  return errors.Trace(err)
			}
		}

		// switch service.Type {
		// case "go":
		// 	containerArgs = append(containerArgs, "rerun", filepath.Join(cnf.Project, service.Workdir, "cmd", arg))
		// }

		if err := watcher.Wait(); err != nil {
		  return errors.Trace(err)
		}

		return nil
	},
}

func resolveDeps(cnf *Config, args []string) ([]string, []string, error) {
	if len(args) == 0 {
		return nil, nil, nil
	}

	services := []string{}
	tools := []string{}
	for _, arg := range args {
		if cnf.IsService(arg) {
			services = append(services, arg)

			subservices, subtools, err := resolveDeps(cnf, cnf.Services[arg].Deps)
			if err != nil {
			  return nil, nil, errors.Trace(err)
			}
			services = append(services, subservices...)
			tools = append(tools, subtools...)
		}

		if cnf.IsTool(arg) {
			tools = append(tools, arg)

			subservices, subtools, err := resolveDeps(cnf, cnf.Tools[arg].Deps)
			if err != nil {
			  return nil, nil, errors.Trace(err)
			}
			services = append(services, subservices...)
			tools = append(tools, subtools...)
		}

		return nil, nil, errors.NotFoundf("service %s", arg)
	}

	services = collections.UniqueStrings(services)
	tools = collections.UniqueStrings(tools)

	return services, tools, nil
}
