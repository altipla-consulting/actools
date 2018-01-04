package main

import (
	"os/exec"
)

func dockerNetworkExists(networkName string) (bool, error) {
	cmd := exec.Command("docker", "network", "inspect", networkName)

	if err := cmd.Start(); err != nil {
		return false, err
	}
	if err := cmd.Wait(); err != nil {
		if !cmd.ProcessState.Success() {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func dockerContainerExists(containerName string) (bool, error) {
	cmd := exec.Command("docker", "inspect", containerName)

	if err := cmd.Start(); err != nil {
		return false, err
	}
	if err := cmd.Wait(); err != nil {
		if !cmd.ProcessState.Success() {
			return false, nil
		}

		return false, err
	}

	return true, nil
}
