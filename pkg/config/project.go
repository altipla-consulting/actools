package config

import (
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

var projectName string

func init() {
	root, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	projectName = filepath.Base(root)
}

func ProjectName() string {
	return projectName
}
