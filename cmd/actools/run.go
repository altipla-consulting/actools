package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func runInteractive(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.WithFields(log.Fields{
		"cmd": append([]string{name}, args...),
	}).Debug("run interactive command")

	if err := cmd.Start(); err != nil {
		return errors.Trace(err)
	}
	if err := cmd.Wait(); err != nil {
		return errors.Trace(err)
	}

	return nil
}

type glideConfig struct {
	Package string `yaml:"package"`
}

func runContainer(container string, args ...string) error {
	container = fmt.Sprintf("eu.gcr.io/altipla-tools/%s:latest", container)

	root, err := os.Getwd()
	if err != nil {
		return errors.Trace(err)
	}

	sh := []string{
		"run", "--rm",
		"-it",
		"--user", fmt.Sprintf("%d:%d", os.Getuid(), os.Getgid()),
		"-v", fmt.Sprintf("%s:/workspace", root),
	}

	sshAuthSock := os.Getenv("SSH_AUTH_SOCK")
	if sshAuthSock == "" {
		log.Warning("WARNING: No SSH_AUTH_SOCK defined in the environment. Start an ssh-agent to share the SSH keys with the tools.")
	}
	if sshAuthSock != "" {
		sh = append(sh, "-e", fmt.Sprintf("SSH_AUTH_SOCK=%s", sshAuthSock))
		sh = append(sh, "-v", fmt.Sprintf("%s:%s", sshAuthSock, sshAuthSock))
	}

	gcloudConfigPath := fmt.Sprintf("%s/.config/gcloud", os.Getenv("HOME"))
	if hasConfig(gcloudConfigPath) {
		sh = append(sh, "-v", fmt.Sprintf("%s:/.config/gcloud", gcloudConfigPath))
	}

	gsutilConfigPath := fmt.Sprintf("%s/.gsutil", os.Getenv("HOME"))
	if hasConfig(gsutilConfigPath) {
		sh = append(sh, "-v", fmt.Sprintf("%s:/.gsutil", gsutilConfigPath))
	}

	kubectlConfigPath := fmt.Sprintf("%s/.kube", os.Getenv("HOME"))
	if hasConfig(kubectlConfigPath) {
		sh = append(sh, "-v", fmt.Sprintf("%s:/.kube", kubectlConfigPath))
	}

	networkName := fmt.Sprintf("%s_default", filepath.Base(root))
	hasNetwork, err := dockerNetworkExists(networkName)
	if err != nil {
		return errors.Trace(err)
	}
	if hasNetwork {
		sh = append(sh, fmt.Sprintf("--network=%s", networkName))
	}

	var project string

	if hasConfig("glide.yaml") {
		f, err := os.Open("glide.yaml")
		if err != nil {
			return errors.Trace(err)
		}

		content, err := ioutil.ReadAll(f)
		if err != nil {
			return errors.Trace(err)
		}

		glideCnf := new(glideConfig)
		if err := yaml.Unmarshal(content, glideCnf); err != nil {
			return errors.Trace(err)
		}

		if glideCnf.Package != "." && glideCnf.Package != "" {
			project = glideCnf.Package
		}
	}

	cnf, err := ReadConfig()
	if err != nil {
	  return errors.Trace(err)
	}
	if cnf != nil {
		project = cnf.Project
	}

	if project != "" {
		sh = append(sh, "-v", fmt.Sprintf("%s:/go/src/%s", root, project))
		sh = append(sh, "-w", fmt.Sprintf("/go/src/%s", project))
	}

	sh = append(sh, container)
	sh = append(sh, args...)

	log.WithFields(log.Fields{
		"sh": sh,
	}).Debug("run docker container")

	if err := runInteractive("docker", sh...); err != nil {
		log.WithFields(log.Fields{"err": err.Error()}).Error("container failed")
		return nil
	}

	return nil
}

func hasConfig(path string) bool {
	if _, err := os.Stat(path); err != nil {
		if !os.IsNotExist(err) {
			fmt.Printf("cannot stat config path <%s>: %s\n", path, err.Error())
			os.Exit(3)
		}

		return false
	}

	return true
}

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
