package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var containers = []string{
	"gcloud",
	"glide",
	"go",
	"gulp",
	"node",
	"protoc",
	"rambler",
	"sass-lint",
}

func init() {
	CmdRoot.AddCommand(CmdPull)
}

var CmdPull = &cobra.Command{
	Use:   "pull",
	Short: "Descarga y actualiza forzosamente las imágenes de los contenedores de herramientas.",
	RunE: func(cmd *cobra.Command, args []string) error {
		for _, container := range containers {
			log.WithFields(log.Fields{"container": container}).Info("download container")
			if err := runInteractive("docker", "pull", fmt.Sprintf("eu.gcr.io/altipla-tools/%s:latest", container)); err != nil {
				return err
			}
		}

		return nil
	},
}
