package containers

import (
	"libs.altipla.consulting/errors"

	"github.com/altipla-consulting/actools/pkg/docker"
)

const Repo = "eu.gcr.io/altipla-tools"

type Container struct {
	Image   string
	Tools   []string
	Options []docker.ContainerOption
}

var containers = []Container{
	{
		Image:   "envoy",
		Tools:   []string{"envoy"},
		Options: []docker.ContainerOption{},
	},
	{
		Image: "cloudsqlproxy",
		Tools: []string{},
		Options: []docker.ContainerOption{
			docker.WithLocalUser(),
			docker.WithSharedGcloud(),
			docker.WithStandardHome(),
		},
	},
	{
		Image: "dev-go",
		Tools: []string{},
		Options: []docker.ContainerOption{
			docker.WithSharedWorkspace(),
			docker.WithLocalUser(),
			docker.WithSharedGopath(),
			docker.WithSharedGcloud(),
			docker.WithStandardHome(),
		},
	},
	{
		Image: "dev-gulp",
		Tools: []string{},
		Options: []docker.ContainerOption{
			docker.WithSharedWorkspace(),
			docker.WithLocalUser(),
			docker.WithStandardHome(),
		},
	},
	{
		Image: "dev-appengine",
		Tools: []string{},
		Options: []docker.ContainerOption{
			docker.WithSharedWorkspace(),
			docker.WithLocalUser(),
			docker.WithSharedSSHSocket(),
			docker.WithSharedGcloud(),
			docker.WithSharedGopath(),
			docker.WithStandardHome(),
		},
	},
	{
		Image: "gcloud",
		Tools: []string{"gcloud", "gsutil", "kubectl"},
		Options: []docker.ContainerOption{
			docker.WithSharedWorkspace(),
			docker.WithLocalUser(),
			docker.WithSharedSSHSocket(),
			docker.WithSharedGcloud(),
			docker.WithStandardHome(),
		},
	},
	{
		Image: "go",
		Tools: []string{"go", "gofmt"},
		Options: []docker.ContainerOption{
			docker.WithSharedWorkspace(),
			docker.WithLocalUser(),
			docker.WithSharedGopath(),
			docker.WithSharedGcloud(),
			docker.WithStandardHome(),
			docker.WithSharedSSHSocket(),
		},
	},
	{
		Image: "gulp",
		Tools: []string{"gulp"},
		Options: []docker.ContainerOption{
			docker.WithSharedWorkspace(),
			docker.WithLocalUser(),
			docker.WithStandardHome(),
		},
	},
	{
		Image: "juice",
		Tools: []string{"juice"},
		Options: []docker.ContainerOption{
			docker.WithSharedWorkspace(),
			docker.WithLocalUser(),
		},
	},
	{
		Image: "mysql",
		Tools: []string{"mysql"},
		Options: []docker.ContainerOption{
			docker.WithSharedWorkspace(),
			docker.WithoutTTY(),
		},
	},
	{
		Image: "mysqldump",
		Tools: []string{"mysqldump"},
		Options: []docker.ContainerOption{
			docker.WithSharedWorkspace(),
			docker.WithLocalUser(),
			docker.WithoutTTY(),
			docker.WithStandardHome(),
		},
	},
	{
		Image: "node",
		Tools: []string{"node", "npm", "npx"},
		Options: []docker.ContainerOption{
			docker.WithSharedWorkspace(),
			docker.WithLocalUser(),
			docker.WithSharedSSHSocket(),
			docker.WithStandardHome(),
		},
	},
	{
		Image:   "phpmyadmin",
		Tools:   []string{},
		Options: []docker.ContainerOption{},
	},
	{
		Image: "protoc",
		Tools: []string{"protoc"},
		Options: []docker.ContainerOption{
			docker.WithSharedWorkspace(),
			docker.WithLocalUser(),
			docker.WithSharedGopath(),
			docker.WithStandardHome(),
		},
	},
	{
		Image: "redis",
		Tools: []string{"redis-cli"},
		Options: []docker.ContainerOption{
			docker.WithoutTTY(),
		},
	},
	{
		Image: "sass-lint",
		Tools: []string{"sass-lint"},
		Options: []docker.ContainerOption{
			docker.WithSharedWorkspace(),
		},
	},
	{
		Image: "migrator",
		Tools: []string{"migrator", "init-migrator"},
		Options: []docker.ContainerOption{
			docker.WithSharedWorkspace(),
			docker.WithStandardHome(),
		},
	},
	{
		Image: "php",
		Tools: []string{"php", "phpunit", "composer"},
		Options: []docker.ContainerOption{
			docker.WithSharedWorkspace(),
		},
	},
	{
		Image:   "prometheus",
		Tools:   []string{},
		Options: []docker.ContainerOption{},
	},
	{
		Image:   "firestore",
		Tools:   []string{},
		Options: []docker.ContainerOption{},
	},
	{
		Image:   "pubsub",
		Tools:   []string{},
		Options: []docker.ContainerOption{},
	},
}

func Images() []string {
	var images []string
	for _, container := range containers {
		images = append(images, container.Image)
	}

	return images
}

func List() []Container {
	return containers
}

func FindContainerTool(tool string) Container {
	for _, container := range containers {
		for _, t := range container.Tools {
			if t == tool {
				return container
			}
		}
	}

	panic("should not reach here")
}

func FindImage(image string) (Container, error) {
	for _, container := range containers {
		if container.Image == image {
			return container, nil
		}
	}

	return Container{}, errors.Errorf("container not found: %s", image)
}
