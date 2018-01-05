package main

import (
	"os/exec"

	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"
)

func dockerNetworkExists(networkName string) (bool, error) {
	cmd := exec.Command("docker", "network", "inspect", networkName)

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

func dockerContainerExists(containerName string) (bool, error) {
	cmd := exec.Command("docker", "inspect", containerName)

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

func dockerCreateNetwork(networkName string) error {
	log.WithFields(log.Fields{"network": networkName}).Info("create network")

	cmd := exec.Command("docker", "network", "create", networkName)
	if err := cmd.Start(); err != nil {
		return errors.Trace(err)
	}
	if err := cmd.Wait(); err != nil {
		return errors.Trace(err)
	}

	return nil
}
