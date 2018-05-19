package main

import (
	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/altipla-consulting/actools/pkg/docker"
)

func init() {
	CmdRoot.AddCommand(CmdStop)
}

var CmdStop = &cobra.Command{
	Use:   "stop",
	Short: "Apaga un servicio de desarrollo. Sin argumentos apaga todos los servicios.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cnf, err := ReadConfig()
		if err != nil {
			return errors.Trace(err)
		}
		if cnf == nil || cnf.Tools == nil {
			return errors.NotFoundf("actools.yml")
		}

		if len(args) == 0 {
			for tool := range cnf.Tools {
				args = append(args, tool)
			}
		}

		for _, arg := range args {
			if _, toolsOk := cnf.Tools[arg]; !toolsOk {
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
