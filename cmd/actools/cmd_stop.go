package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	CmdRoot.AddCommand(CmdStop)
}

var CmdStop = &cobra.Command{
	Use:   "stop",
	Short: "Apaga un servicio de desarrollo",
	RunE: func(cmd *cobra.Command, args []string) error {
		cnf, err := ReadConfig()
		if err != nil {
			return errors.Trace(err)
		}
		if cnf == nil {
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
		}

		for _, arg := range args {
			if _, ok := cnf.Tools[arg]; !ok {
				return errors.NotFoundf("tool %s", arg)
			}
		}

		for _, arg := range args {
			name := fmt.Sprintf("%s_%s", root, arg)

			hasContainer, err := dockerContainerExists(name)
			if err != nil {
				return errors.Trace(err)
			}
			if !hasContainer {
				log.WithFields(log.Fields{"tool": arg}).Info("tool already removed")
				return nil
			}

			log.WithFields(log.Fields{"tool": arg}).Info("stop tool")
			if err := runInteractiveDebugOutput("docker", "stop", name); err != nil {
				return errors.Trace(err)
			}
		}

		return nil
	},
}
