package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"libs.altipla.consulting/errors"

	"github.com/altipla-consulting/actools/pkg/containers"
	"github.com/altipla-consulting/actools/pkg/docker"
	"github.com/altipla-consulting/actools/pkg/run"
)

const GoVersion = "go1.13.8"

func init() {
	CmdRoot.AddCommand(CmdPull)
}

var CmdPull = &cobra.Command{
	Use:   "pull",
	Short: "Descarga y actualiza forzosamente las imágenes de los contenedores de herramientas.",
	RunE: func(cmd *cobra.Command, args []string) error {
		var ver string
		output, err := run.InteractiveCaptureOutput("go", "version")
		if err != nil {
			if strings.HasPrefix(errors.Cause(err).Error(), `exec: "go": executable file not found in $PATH`) {
				ver = "NOT_FOUND"
			} else {
				return errors.Trace(err)
			}
		} else {
			ver = strings.Split(output, " ")[2]
		}
		if ver != GoVersion {
			log.WithFields(log.Fields{
				"installed": ver,
				"want":      GoVersion,
			}).Info("Found different version of Go, installing the latest one")

			log.Info("Remove old installation if present")
			if err := run.Interactive("sudo", "rm", "-rf", "/usr/local/go"); err != nil {
				return errors.Trace(err)
			}

			log.Info("Download compressed file")
			if err := run.NonInteractiveWithOutput("wget", "-q", "-O", "/tmp/go.tar.gz", "https://dl.google.com/go/"+GoVersion+".linux-amd64.tar.gz"); err != nil {
				return errors.Trace(err)
			}

			log.Info("Extract file")
			if err := run.Interactive("sudo", "tar", "-C", "/usr/local", "-xzf", "/tmp/go.tar.gz"); err != nil {
				return errors.Trace(err)
			}

			if os.Getenv("GOROOT") == "" {
				bashrc := filepath.Join(os.Getenv("HOME"), ".bashrc")
				log.WithField("path", bashrc).Info("Save info in .bashrc")
				content, err := ioutil.ReadFile(bashrc)
				if err != nil {
					return errors.Trace(err)
				}

				buf := bytes.NewBuffer(content)
				fmt.Fprintln(buf)
				fmt.Fprintln(buf, "export GOROOT=/usr/local/go")
				fmt.Fprintln(buf, "export PATH=$PATH:$GOROOT/bin")

				if err := ioutil.WriteFile(bashrc, buf.Bytes(), 0600); err != nil {
					return errors.Trace(err)
				}

				log.Warning()
				log.Warning("Run the following command or restart the shell:")
				log.Warning()
				log.Warning("\tsource ~/.bashrc")
				log.Warning()
			}
		}

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
