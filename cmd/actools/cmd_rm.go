package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"libs.altipla.consulting/errors"

	"github.com/altipla-consulting/actools/pkg/config"
	"github.com/altipla-consulting/actools/pkg/docker"
)

func init() {
	CmdRoot.AddCommand(CmdRm)
}

var CmdRm = &cobra.Command{
	Use:   "rm",
	Short: "Elimina un servicio de desarrollo. Sin argumentos elimina todos los servicios.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			for tool := range config.Settings.Tools {
				args = append(args, tool)
			}
			for service := range config.Settings.Services {
				args = append(args, service)
			}
		}

		for _, arg := range args {
			if !config.Settings.IsService(arg) && !config.Settings.IsTool(arg) {
				return errors.Errorf("service not found: %s", arg)
			}
		}

		for _, service := range args {
			container, err := docker.Container(service)
			if err != nil {
				return errors.Trace(err)
			}

			if exists, err := container.Exists(); err != nil {
				return errors.Trace(err)
			} else if !exists {
				continue
			}

			if running, err := container.Running(); err != nil {
				return errors.Trace(err)
			} else if running {
				log.WithField("service", service).Info("Stop service")
				if err := container.Stop(); err != nil {
					return errors.Trace(err)
				}
			}

			log.WithField("service", service).Info("Remove service")
			if err := container.Remove(); err != nil {
				return errors.Trace(err)
			}
		}

		return nil
	},
}
