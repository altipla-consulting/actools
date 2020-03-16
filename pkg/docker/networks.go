package docker

import (
	"os/exec"

	log "github.com/sirupsen/logrus"
	"libs.altipla.consulting/errors"
)

type NetworkManager struct {
	name string
}

func Network(name string) *NetworkManager {
	return &NetworkManager{name: name}
}

func (network *NetworkManager) Exists() (bool, error) {
	cmd := exec.Command("docker", "network", "inspect", network.name)
	if err := cmd.Run(); err != nil {
		if !cmd.ProcessState.Success() {
			return false, nil
		}

		return false, errors.Trace(err)
	}

	return true, nil
}

func (network *NetworkManager) Create() error {
	log.WithFields(log.Fields{"network": network.name}).Info("Create network")

	cmd := exec.Command("docker", "network", "create", network.name)
	return errors.Trace(cmd.Run())
}

func (network *NetworkManager) CreateIfNotExists() error {
	exists, err := network.Exists()
	if err != nil {
		return errors.Trace(err)
	}

	if !exists {
		return errors.Trace(network.Create())
	}

	return nil
}

func (network *NetworkManager) String() string {
	return network.name
}
