package main

import (
	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/altipla-consulting/actools/pkg/config"
	"github.com/altipla-consulting/actools/pkg/docker"
)

func init() {
	CmdRoot.AddCommand(CmdStop)
}

var CmdStop = &cobra.Command{
	Use:   "stop",
	Short: "Apaga un servicio de desarrollo. Sin argumentos apaga todos los servicios.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			for tool := range config.Settings.Tools {
				args = append(args, tool)
			}
		}

		for _, arg := range args {
			if !config.Settings.IsTool(arg) {
				return errors.NotFoundf("service %s", arg)
			}
		}

		for _, service := range args {
			container, err := docker.Container(service)
			if err != nil {
				return errors.Trace(err)
			}

			if running, err := container.Running(); err != nil {
				return errors.Trace(err)
			} else if !running {
				continue
			}

			log.WithField("service", service).Info("Stop service")
			if err := container.Stop(); err != nil {
				return errors.Trace(err)
			}
		}

		return nil
	},
}
