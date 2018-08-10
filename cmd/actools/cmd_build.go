package main

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/altipla-consulting/actools/pkg/config"
	"github.com/altipla-consulting/actools/pkg/docker"
	"github.com/altipla-consulting/actools/pkg/run"
)

func init() {
	CmdRoot.AddCommand(CmdBuild)
}

type rootSettings struct {
	Build   buildSettings   `yaml:"build"`
	Release releaseSettings `yaml:"release"`
}

func (settings rootSettings) Valid() bool {
	return settings.Build.Dockerfile != "" && settings.Build.Repo != "" && settings.Build.Project != "" && settings.Build.Name != "" && settings.Release.Deployment != ""
}

type buildSettings struct {
	Dockerfile string `yaml:"dockerfile"`
	Repo       string `yaml:"repo"`
	Project    string `yaml:"project"`
	Name       string `yaml:"name"`
}

type releaseSettings struct {
	Deployment string `yaml:"deployment"`
	Service    string `yaml:"service"`
	ConfigMap  string `yaml:"configmap"`
}

var CmdBuild = &cobra.Command{
	Use:   "build",
	Short: "Construye una imagen y la despliega en producción si ha cambiado.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		content, err := ioutil.ReadFile(args[0])
		if err != nil {
			return errors.Trace(err)
		}

		settings := rootSettings{}
		if err := yaml.Unmarshal(content, &settings); err != nil {
			return errors.Trace(err)
		}

		if !settings.Valid() {
			return errors.NotValidf("settings")
		}

		container := fmt.Sprintf("%s/%s", settings.Build.Project, settings.Build.Name)
		image := docker.Image(settings.Build.Repo, container)
		if err := image.Build(".", settings.Build.Dockerfile); err != nil {
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

		cm, err := getConfigMapHash(settings)
		if err != nil {
			return errors.Trace(err)
		}

		cacheTag := fmt.Sprintf("%s::%s", version, cm)
		if cache.Containers[container] != cacheTag {
			log.WithFields(log.Fields{
				"cached":    cache.Containers[container],
				"new":       version,
				"configmap": cm,
			}).Info("Version changed, pushing to production")

			if err := image.Push(version); err != nil {
				return errors.Trace(err)
			}
			if err := image.Push("latest"); err != nil {
				return errors.Trace(err)
			}

			content, err := ioutil.ReadFile(settings.Release.Deployment)
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
			content = bytes.Replace(content, []byte("cmversion"), []byte(cm), -1)

			if _, err := f.Write(content); err != nil {
				return errors.Trace(err)
			}

			if err := run.InteractiveWithOutput("kubectl", "apply", "-f", f.Name()); err != nil {
				return errors.Trace(err)
			}

			if settings.Release.Service != "" {
				if err := run.InteractiveWithOutput("kubectl", "apply", "-f", settings.Release.Service); err != nil {
					return errors.Trace(err)
				}
			}

			cache.Containers[container] = cacheTag
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

func getConfigMapHash(settings rootSettings) (string, error) {
	if settings.Release.ConfigMap == "" {
		return "", nil
	}

	f, err := os.Open(settings.Release.ConfigMap)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
