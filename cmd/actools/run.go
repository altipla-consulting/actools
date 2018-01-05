package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"github.com/altipla-consulting/actools/pkg/docker"
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

func runInteractiveDebugOutput(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin

	loggerOut := log.New()
	loggerOut.SetLevel(log.StandardLogger().Level)
	wout := loggerOut.WriterLevel(log.DebugLevel)
	defer wout.Close()
	cmd.Stdout = wout

	loggerErr := log.New()
	loggerErr.SetLevel(log.StandardLogger().Level)
	werr := loggerErr.WriterLevel(log.ErrorLevel)
	defer werr.Close()
	cmd.Stderr = werr

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

type containerConfig struct {
	ShareWorkspace    bool
	LocalUser         bool
	ShareSSHSocket    bool
	ShareGcloudConfig bool
	ConfigureGopath   bool
	Persistent        bool
	Ports             []string
	Volumes           []string
	Name              string
	CreateOnly        bool
	NetworkAlias      []string
}

func runContainer(container string, cnf *containerConfig, args ...string) error {
	if cnf == nil {
		cnf = new(containerConfig)
	}

	root, err := os.Getwd()
	if err != nil {
		return errors.Trace(err)
	}

	sh := []string{}

	if cnf.CreateOnly {
		sh = append(sh, "create")
	} else {
		sh = append(sh, "run")
	}

	sh = append(sh, "-it")

	if cnf.Name != "" {
		sh = append(sh, "--name", cnf.Name)
	}

	if !cnf.Persistent {
		sh = append(sh, "--rm")
	}

	if cnf.ShareWorkspace {
		sh = append(sh, "-v", fmt.Sprintf("%s:/workspace", root))
	}

	if cnf.LocalUser {
		sh = append(sh, "--user", fmt.Sprintf("%d:%d", os.Getuid(), os.Getgid()))
	}

	for _, alias := range cnf.NetworkAlias {
		sh = append(sh, "--network-alias", alias)
	}

	if cnf.ShareSSHSocket {
		sshAuthSock := os.Getenv("SSH_AUTH_SOCK")
		if sshAuthSock == "" {
			log.Warning("WARNING: No SSH_AUTH_SOCK defined in the environment. Start an ssh-agent to share the SSH keys with the tools.")
		}
		if sshAuthSock != "" {
			sh = append(sh, "-e", fmt.Sprintf("SSH_AUTH_SOCK=%s", sshAuthSock))
			sh = append(sh, "-v", fmt.Sprintf("%s:%s", sshAuthSock, sshAuthSock))
		}
	}

	if cnf.ShareGcloudConfig {
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
	}

	if cnf.ConfigureGopath {
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
			hostBin := fmt.Sprintf("/tmp/actools-cache-%s/bin", filepath.Base(root))
			sh = append(sh, "-v", fmt.Sprintf("%s:/go/bin", hostBin))
			if err := os.MkdirAll(hostBin, 0777); err != nil {
				return errors.Trace(err)
			}

			hostPkg := fmt.Sprintf("/tmp/actools-cache-%s/pkg", filepath.Base(root))
			sh = append(sh, "-v", fmt.Sprintf("%s:/go/pkg", hostPkg))
			if err := os.MkdirAll(hostPkg, 0777); err != nil {
				return errors.Trace(err)
			}

			sh = append(sh, "-v", fmt.Sprintf("%s:/go/src/%s", root, project))
			sh = append(sh, "-w", fmt.Sprintf("/go/src/%s", project))
		}
	}

	networkName := fmt.Sprintf("%s_default", filepath.Base(root))
	hasNetwork, err := docker.NetworkExists(networkName)
	if err != nil {
		return errors.Trace(err)
	}
	if hasNetwork {
		sh = append(sh, fmt.Sprintf("--network=%s", networkName))
	}

	for _, port := range cnf.Ports {
		sh = append(sh, "-p", port)
	}
	for _, volume := range cnf.Volumes {
		if strings.HasPrefix(volume, "./") {
			volume = filepath.Join(root, volume[2:])
		}
		sh = append(sh, "-v", volume)
	}
	sh = append(sh, fmt.Sprintf("eu.gcr.io/altipla-tools/%s:latest", container))
	sh = append(sh, args...)

	log.WithFields(log.Fields{"sh": sh}).Debug("run docker container")

	if cnf.CreateOnly {
		if err := runInteractiveDebugOutput("docker", sh...); err != nil {
			log.WithFields(log.Fields{"err": err.Error()}).Error("container failed")
			return nil
		}
	} else {
		if err := runInteractive("docker", sh...); err != nil {
			log.WithFields(log.Fields{"err": err.Error()}).Error("container failed")
			return nil
		}
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
