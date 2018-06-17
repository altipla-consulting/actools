package docker

import (
	"fmt"
	"path"
	"strings"

	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"

	"github.com/altipla-consulting/actools/pkg/run"
)

type ImageManager struct {
	name string
}

func Image(repo, name string) *ImageManager {
	return &ImageManager{
		name: fmt.Sprintf("%s/%s", repo, name),
	}
}

func (image *ImageManager) Pull() error {
	return errors.Trace(run.InteractiveWithOutput("docker", "pull", image.String()))
}

func (image *ImageManager) Push(tag string) error {
	log.WithFields(log.Fields{
		"name": image.name,
		"tag":  tag,
	}).Info("Push image")

	taggedName := fmt.Sprintf("%s:%s", image.name, tag)
	if err := run.InteractiveWithOutput("docker", "tag", image.name, taggedName); err != nil {
		return errors.Trace(err)
	}

	return errors.Trace(run.InteractiveWithOutput("docker", "push", taggedName))
}

func (image *ImageManager) Build(context, dockerfile string) error {
	log.WithFields(log.Fields{
		"context":    context,
		"dockerfile": dockerfile,
		"name":       image.name,
	}).Info("Build image")

	if dockerfile == "" {
		dockerfile = path.Join(context, "Dockerfile")
	}

	return errors.Trace(run.InteractiveWithOutput("docker", "build", "-t", image.name, "-f", dockerfile, context))
}

func (image *ImageManager) LastBuiltID() (string, error) {
	version, err := run.InteractiveCaptureOutput("docker", "image", "inspect", image.name, "-f", "{{.Id}}")
	if err != nil {
		return "", errors.Trace(err)
	}

	return strings.Split(version, ":")[1][:12], nil
}

func (image *ImageManager) String() string {
	return fmt.Sprintf("%s:latest", image.name)
}
