package docker

import (
	"os/exec"

	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"
)

func NetworkExists(name string) (bool, error) {
	cmd := exec.Command("docker", "network", "inspect", name)

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

func CreateNetwork(name string) error {
	log.WithFields(log.Fields{"network": name}).Info("Create network")

	cmd := exec.Command("docker", "network", "create", name)
	if err := cmd.Start(); err != nil {
		return errors.Trace(err)
	}
	if err := cmd.Wait(); err != nil {
		return errors.Trace(err)
	}

	return nil
}
