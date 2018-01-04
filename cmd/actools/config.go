package main

import (
	"os"
	"io/ioutil"

	"github.com/juju/errors"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Project string `yaml:"project"`

	Services map[string]*Service `yaml:"services"`
	Tools map[string]*Tool `yaml:"tool"`
}

type Service struct {
	Lang string `yaml:"lang"`
	Deps []string `yaml:"deps"`
	Ports []string `yaml:"ports"`
}

type Tool struct {
	Container string `yaml:"container"`
	Ports []string `yaml:"ports"`
}

func ReadConfig() (*Config, error) {
	f, err := os.Open("actools.yaml")
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
