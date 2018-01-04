package main

import (
	"io/ioutil"
	"os"

	"github.com/juju/errors"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Project string `yaml:"project"`

	Services map[string]*Service `yaml:"services"`
}

type Service struct {
	Container  string   `yaml:"container"`
	Deps       []string `yaml:"deps"`
	Ports      []string `yaml:"ports"`
	Volumes    []string `yaml:"volumes"`
	Persistent bool     `yaml:"persistent"`
}

func ReadConfig() (*Config, error) {
	f, err := os.Open("actools.yml")
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}

		return nil, errors.Trace(err)
	}

	content, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, errors.Trace(err)
	}

	cnf := new(Config)
	if err := yaml.Unmarshal(content, cnf); err != nil {
		return nil, errors.Trace(err)
	}

	return cnf, nil
}
