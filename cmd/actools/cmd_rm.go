package main

import (
	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/altipla-consulting/actools/pkg/docker"
)

func init() {
	CmdRoot.AddCommand(CmdRm)
}

var CmdRm = &cobra.Command{
	Use:   "rm",
	Short: "Elimina un servicio de desarrollo persistente. Sin argumentos elimina todos los servicios.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cnf, err := ReadConfig()
		if err != nil {
			return errors.Trace(err)
		}
		if cnf == nil || (cnf.Services == nil && cnf.Tools == nil) {
			return errors.NotFoundf("actools.yml")
		}

		if len(args) == 0 {
			for tool := range cnf.Tools {
				args = append(args, tool)
			}
			for service := range cnf.Services {
				args = append(args, service)
			}
		}

		for _, arg := range args {
			_, serviceOk := cnf.Services[arg]
			_, toolsOk := cnf.Tools[arg]
			if !serviceOk && !toolsOk {
				return errors.NotFoundf("service %s", arg)
			}
		}

		for _, service := range args {
			container, err := docker.Container(service)
			if err != nil {
				return errors.Trace(err)
			}

			exists, err := container.Exists()
			if err != nil {
				return errors.Trace(err)
			}
			if !exists {
				log.WithField("service", service).Info("Remove service")
				continue
			}

			log.WithField("service", service).Info("Stop service")
			if err := container.Stop(); err != nil {
				return errors.Trace(err)
			}

			log.WithField("service", service).Info("Remove service")
			if err := container.Remove(); err != nil {
				return errors.Trace(err)
			}
		}

		return nil
	},
}
