package docker

import (
	"os/exec"

	"github.com/juju/errors"
)

func ContainerExists(name string) (bool, error) {
	cmd := exec.Command("docker", "inspect", name)

	if err := cmd.Start(); err != nil {
		return false, errors.Trace(err)
	}
	if err := cmd.Wait(); err != nil {
		if !cmd.ProcessState.Success() {
			return false, nil
		}

		return false, errors.Trace(err)
	}

	return true, nil
}
