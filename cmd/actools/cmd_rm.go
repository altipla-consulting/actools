package main

import (
	"fmt"
	"os"
	"path/filepath"

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
	Short: "Elimina un servicio de desarrollo persistente",
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
		root = filepath.Base(root)

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
				return errors.NotFoundf("tool %s", arg)
			}
		}

		for _, arg := range args {
			name := fmt.Sprintf("%s_%s", root, arg)

			hasContainer, err := docker.ContainerExists(name)
			if err != nil {
				return errors.Trace(err)
			}
			if !hasContainer {
				log.WithFields(log.Fields{"tool": arg}).Info("Tool already removed")
				continue
			}

			log.WithFields(log.Fields{"tool": arg}).Info("Stop tool")
			if err := runInteractiveDebugOutput("docker", "stop", name); err != nil {
				return errors.Trace(err)
			}

			log.WithFields(log.Fields{"tool": arg}).Info("Remove tool")
			if err := runInteractiveDebugOutput("docker", "rm", name); err != nil {
				return errors.Trace(err)
			}
		}

		return nil
	},
}
