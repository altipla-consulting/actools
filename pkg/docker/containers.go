package docker

import (
	"fmt"
	"os/exec"

	"github.com/juju/errors"

	"github.com/altipla-consulting/actools/pkg/config"
	"github.com/altipla-consulting/actools/pkg/run"
)

type ContainerManager struct {
	name    string
	network *NetworkManager
}

func Container(name string, options ...ContainerOption) *ContainerManager {
	container := &ContainerManager{
		name: fmt.Sprintf("%s_%s", config.ProjectName(), name),
	}

	for _, option := range options {
		option(container)
	}

	return container
}

func (container *ContainerManager) Exists() (bool, error) {
	cmd := exec.Command("docker", "inspect", container.name)
	if err := cmd.Run(); err != nil {
		if !cmd.ProcessState.Success() {
			return false, nil
		}

		return false, errors.Trace(err)
	}

	return true, nil
}

func (container *ContainerManager) Stop() error {
	exists, err := container.Exists()
	if err != nil {
		return errors.Trace(err)
	}
	if !exists {
		return nil
	}

	return errors.Trace(run.Interactive("docker", "stop", container.name))
}

func (container *ContainerManager) Kill() error {
	exists, err := container.Exists()
	if err != nil {
		return errors.Trace(err)
	}
	if !exists {
		return nil
	}

	return errors.Trace(run.Interactive("docker", "kill", container.name))
}

func (container *ContainerManager) Remove() error {
	exists, err := container.Exists()
	if err != nil {
		return errors.Trace(err)
	}
	if !exists {
		return nil
	}

	return errors.Trace(run.Interactive("docker", "rm", container.name))
}
