package config

import (
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

var (
	projectName    string
	projectPackage string
)

func init() {
	if Settings.Project != "" {
		projectPackage = Settings.Project
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
