package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"libs.altipla.consulting/errors"

	"github.com/altipla-consulting/actools/pkg/containers"
	"github.com/altipla-consulting/actools/pkg/docker"
)

const GoVersion = "go1.13.8"

func init() {
	CmdRoot.AddCommand(CmdPull)
}

var CmdPull = &cobra.Command{
	Use:   "pull",
	Short: "Descarga y actualiza forzosamente las imágenes de los contenedores de herramientas.",
	RunE: func(cmd *cobra.Command, args []string) error {
		for _, image := range containers.Images() {
			log.WithField("image", image).Info("Download image")

			image := docker.Image("eu.gcr.io", fmt.Sprintf("altipla-tools/%s", image))
			if err := image.Pull(); err != nil {
				return errors.Trace(err)
			}
		}

		return nil
	},
}
