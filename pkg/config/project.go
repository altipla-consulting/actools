package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

var (
	projectName    string
	projectPackage string
)

type glideConfig struct {
	Package string `yaml:"package"`
}

func init() {
	// First source: the glide.yaml file with its package name
	content, err := ioutil.ReadFile("glide.yaml")
	if err != nil {
		log.Fatal(err)
	}
	glideCnf := new(glideConfig)
	if err := yaml.Unmarshal(content, glideCnf); err != nil {
		log.Fatal(err)
	}
	if glideCnf.Package != "." && glideCnf.Package != "" {
		projectPackage = glideCnf.Package
	}

	cnf, err := ReadConfig()
	if err != nil {
		log.Fatal(err)
	}
	if cnf != nil {
		projectPackage = cnf.Project
	}

	root, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	projectName = filepath.Base(root)
}

func ProjectName() string {
	return projectName
}

func ProjectPackage() string {
	return projectPackage
}
