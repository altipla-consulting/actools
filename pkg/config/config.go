package config

import (
	"io/ioutil"
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var Settings = new(Config)

func init() {
	content, err := ioutil.ReadFile("actools.yml")
	if err != nil {
		if os.IsNotExist(err) {
			return
		}

		log.Fatal(err)
	}

	if err := yaml.Unmarshal(content, &Settings); err != nil {
		log.Fatal(err)
	}
}

type Config struct {
	Project string `yaml:"project"`

	Services map[string]*Service `yaml:"services"`
	Tools    map[string]*Tool    `yaml:"tools"`
}

func (cnf *Config) IsService(name string) bool {
	_, ok := cnf.Services[name]
	return ok
}

func (cnf *Config) IsTool(name string) bool {
	_, ok := cnf.Tools[name]
	return ok
}

type Service struct {
	Type    string            `yaml:"type"`
	Deps    []string          `yaml:"deps"`
	Ports   []string          `yaml:"ports"`
	Workdir string            `yaml:"workdir"`
	Volumes []string          `yaml:"volumes"`
	Env     map[string]string `yaml:"env"`
	Ignore  []string          `yaml:"ignore"`
}

type Tool struct {
	Container string   `yaml:"container"`
	Deps      []string `yaml:"deps"`
	Ports     []string `yaml:"ports"`
	Volumes   []string `yaml:"volumes"`
	Args      []string `yaml:"args"`
}
