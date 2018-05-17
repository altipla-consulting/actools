package containers

import (
	"github.com/altipla-consulting/actools/pkg/docker"
)

type Container struct {
	Image   string
	Tools   []string
	Options []docker.ContainerOption
}

var containers = []Container{
	{
		Image:   "baster",
		Tools:   []string{"baster"},
		Options: []docker.ContainerOption{},
	},
	{
		Image:   "cloudsqlproxy",
		Tools:   []string{},
		Options: []docker.ContainerOption{},
	},
	{
		Image: "dev-go",
		Tools: []string{},
		Options: []docker.ContainerOption{
			docker.WithSharedWorkspace(),
			docker.WithLocalUser(),
			docker.WithSharedGopath(),
			docker.WithSharedGcloud(),
		},
	},
	{
		Image:   "dev-gulp",
		Tools:   []string{},
		Options: []docker.ContainerOption{},
	},
	{
		Image: "gcloud",
		Tools: []string{"gcloud", "gsutil", "kubectl"},
		Options: []docker.ContainerOption{
			docker.WithSharedWorkspace(),
			docker.WithLocalUser(),
			docker.WithSharedSSHSocket(),
			docker.WithSharedGcloud(),
		},
	},
	{
		Image: "glide",
		Tools: []string{"glide"},
		Options: []docker.ContainerOption{
			docker.WithSharedWorkspace(),
			docker.WithLocalUser(),
			docker.WithSharedSSHSocket(),
			docker.WithSharedGopath(),
		},
	},
	{
		Image: "go",
		Tools: []string{"go", "gofmt"},
		Options: []docker.ContainerOption{
			docker.WithSharedWorkspace(),
			docker.WithLocalUser(),
			docker.WithSharedGopath(),
		},
	},
	{
		Image: "gulp",
		Tools: []string{"gulp"},
		Options: []docker.ContainerOption{
			docker.WithSharedWorkspace(),
			docker.WithLocalUser(),
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
		Image: "king",
		Tools: []string{"king"},
		Options: []docker.ContainerOption{
			docker.WithSharedWorkspace(),
			docker.WithLocalUser(),
			docker.WithSharedGopath(),
		},
	},
	{
		Image: "mysql",
		Tools: []string{"mysql", "mysqldump"},
		Options: []docker.ContainerOption{
			docker.WithSharedWorkspace(),
			docker.WithLocalUser(),
			docker.WithoutTTY(),
		},
	},
	{
		Image: "node",
		Tools: []string{"node", "npm"},
		Options: []docker.ContainerOption{
			docker.WithSharedWorkspace(),
			docker.WithLocalUser(),
			docker.WithSharedSSHSocket(),
		},
	},
	{
		Image:   "pdfgen",
		Tools:   []string{"pdfgen"},
		Options: []docker.ContainerOption{},
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
		},
	},
	{
		Image: "rambler",
		Tools: []string{"rambler"},
		Options: []docker.ContainerOption{
			docker.WithSharedWorkspace(),
		},
	},
	{
		Image:   "redis",
		Tools:   []string{},
		Options: []docker.ContainerOption{},
	},
	{
		Image: "sass-lint",
		Tools: []string{"sass-lint"},
		Options: []docker.ContainerOption{
			docker.WithSharedWorkspace(),
		},
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
