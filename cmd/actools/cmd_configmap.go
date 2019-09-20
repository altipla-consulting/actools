package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"libs.altipla.consulting/errors"
)

var configMapName string

func init() {
	CmdConfigMap.PersistentFlags().StringVar(&configMapName, "name", "", "Nombre del ConfigMap que vamos a generar")
	CmdRoot.AddCommand(CmdConfigMap)
}

var CmdConfigMap = &cobra.Command{
	Use:   "configmap",
	Short: "Genera un ConfigMap de Kubernetes a partir de uno o varios ficheros",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return errors.New("empty args")
		}
		if configMapName == "" {
			return errors.New("--name argument is required")
		}

		cnf := &cmConfig{
			Kind:       "ConfigMap",
			APIVersion: "v1",
			Data:       make(map[string]string),
			Metadata:   metadata{Name: configMapName},
		}
		for _, arg := range args {
			content, err := ioutil.ReadFile(arg)
			if err != nil {
				return errors.Trace(err)
			}
			cnf.Data[filepath.Base(arg)] = string(content)
		}

		output, err := yaml.Marshal(cnf)
		if err != nil {
			return errors.Trace(err)
		}

		fmt.Println(string(output))

		return nil
	},
}

type cmConfig struct {
	APIVersion string            `yaml:"apiVersion"`
	Kind       string            `yaml:"kind"`
	Metadata   metadata          `yaml:"metadata"`
	Data       map[string]string `yaml:"data"`
}

type metadata struct {
	Name string `yaml:"name"`
}
