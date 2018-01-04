package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var foreground bool

func init() {
	CmdRoot.AddCommand(CmdStart)
	CmdStart.PersistentFlags().BoolVar(&foreground, "fg", false, "Ejecuta el contenedor en primer plano para ver los logs")
}

var CmdStart = &cobra.Command{
	Use:   "start",
	Short: "Enciende un servicio de desarrollo",
	RunE: func(cmd *cobra.Command, args []string) error {
		cnf, err := ReadConfig()
		if err != nil {
			return errors.Trace(err)
		}
		if cnf == nil || cnf.Services == nil {
			return errors.NotFoundf("actools.yml")
		}

		root, err := os.Getwd()
		if err != nil {
			return errors.Trace(err)
		}

		for _, arg := range args {
			if _, ok := cnf.Services[arg]; !ok {
				return errors.NotFoundf("service %s", arg)
			}
		}

		for _, arg := range args {
			service := cnf.Services[arg]

			name := fmt.Sprintf("%s_%s", filepath.Base(root), arg)
			hasContainer, err := dockerContainerExists(name)
			if err != nil {
				return errors.Trace(err)
			}
			log.WithFields(log.Fields{
				"service":    arg,
				"persistent": service.Persistent,
				"is-new":     !hasContainer,
			}).Info("start service")

			if hasContainer {
				if err := runInteractive("docker", "start", name); err != nil {
					return errors.Trace(err)
				}
				continue
			}

			containerArgs := []string{
				"--name", name,
			}
			if !foreground {
				containerArgs = append(containerArgs, "-d")
			}
			for _, port := range service.Ports {
				containerArgs = append(containerArgs, "-p", port)
			}
			for _, volume := range service.Volumes {
				if strings.HasPrefix(volume, "./") {
					volume = filepath.Join(root, volume[2:])
				}
				containerArgs = append(containerArgs, "-v", volume)
			}

			containerCnf := &containerConfig{
				DockerArgs: containerArgs,
				Persistent: service.Persistent,
			}
			if err := runContainer(cnf.Services[arg].Container, containerCnf); err != nil {
				return errors.Annotatef(err, "service %s", arg)
			}
		}

		return nil
	},
}
