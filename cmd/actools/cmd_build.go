package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/altipla-consulting/actools/pkg/config"
	"github.com/altipla-consulting/actools/pkg/docker"
	"github.com/altipla-consulting/actools/pkg/run"
)

var k8sfile, dockerfile string
var repo, container string

func init() {
	CmdRoot.AddCommand(CmdBuild)
	CmdBuild.PersistentFlags().StringVarP(&k8sfile, "kubernetes", "", "", "Deployment de Kubernetes para el despliegue")
	CmdBuild.PersistentFlags().StringVarP(&dockerfile, "dockerfile", "", "", "Dockerfile para la construcción")
	CmdBuild.PersistentFlags().StringVarP(&repo, "repo", "", "", "Repositorio para el despliegue")
	CmdBuild.PersistentFlags().StringVarP(&container, "name", "", "", "Nombre de la imagen")
}

var CmdBuild = &cobra.Command{
	Use:   "build",
	Short: "Construye una imagen y la despliega en producción si ha cambiado.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if k8sfile == "" || dockerfile == "" || repo == "" || container == "" {
			return errors.NotValidf("arguments required")
		}

		image := docker.Image(repo, container)
		if err := image.Build(".", dockerfile); err != nil {
			return errors.Trace(err)
		}

		version, err := image.LastBuiltID()
		if err != nil {
			return errors.Trace(err)
		}

		cache, err := getBuildCache()
		if err != nil {
			return errors.Trace(err)
		}

		if cache.Containers[container] != version {
			log.WithFields(log.Fields{
				"cached": cache.Containers[container],
				"new":    version,
			}).Info("Version changed, pushing to production")

			if err := image.Push(version); err != nil {
				return errors.Trace(err)
			}
			if err := image.Push("latest"); err != nil {
				return errors.Trace(err)
			}

			content, err := ioutil.ReadFile(k8sfile)
			if err != nil {
				return errors.Trace(err)
			}
			f, err := ioutil.TempFile("", "deployment.yaml")
			if err != nil {
				return errors.Trace(err)
			}
			defer f.Close()
			defer os.Remove(f.Name())

			content = bytes.Replace(content, []byte("latest"), []byte(version), -1)

			if _, err := f.Write(content); err != nil {
				return errors.Trace(err)
			}

			if err := run.InteractiveWithOutput("kubectl", "apply", "-f", f.Name()); err != nil {
				return errors.Trace(err)
			}

			cache.Containers[container] = version
			if err := cache.Save(); err != nil {
				return errors.Trace(err)
			}
		}

		return nil
	},
}

type buildCache struct {
	Containers map[string]string
}

func (cache *buildCache) Save() error {
	filename := filepath.Join(config.Home(), ".actools", "build-cache.json")
	f, err := os.Create(filename)
	if err != nil {
		return errors.Trace(err)
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(cache); err != nil {
		return errors.Trace(err)
	}

	return nil
}

func getBuildCache() (*buildCache, error) {
	filename := filepath.Join(config.Home(), ".actools", "build-cache.json")
	f, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return &buildCache{
				Containers: map[string]string{},
			}, nil
		}

		return nil, errors.Trace(err)
	}
	defer f.Close()

	cache := new(buildCache)
	if err := json.NewDecoder(f).Decode(&cache); err != nil {
		return nil, errors.Trace(err)
	}

	return cache, nil
}
