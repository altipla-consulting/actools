package docker

import (
	"fmt"

	"github.com/juju/errors"

	"github.com/altipla-consulting/actools/pkg/run"
)

type ImageManager struct {
	name string
}

func Image(repo, name, tag string) *ImageManager {
	return &ImageManager{
		name: fmt.Sprintf("%s/%s:%s", repo, name, tag),
	}
}

func (image *ImageManager) Pull() error {
	return errors.Trace(run.InteractiveWithOutput("docker", "pull", image.name))
}

func (image *ImageManager) String() string {
	return image.name
}
