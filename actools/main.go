package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/altipla-consulting/collections"
	"gopkg.in/yaml.v2"
)

var containers = []string{
	"gcloud",
	"glide",
	"go",
	"gulp",
	"node",
	"protoc",
	"rambler",
	"sass-lint",
}

var toolContainers = map[string]string{
	"npm":     "node",
	"gsutil":  "gcloud",
	"kubectl": "gcloud",
	"gofmt":   "go",
}

type glideConfig struct {
	Package string `yaml:"package"`
}

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Pass the tool name as the first argument: actools glide")
		os.Exit(1)
	}

	containerName, ok := toolContainers[os.Args[1]]
	if !ok {
		containerName = os.Args[1]
	}
	if containerName == "debug" {
		containerName = os.Args[2]
	}
	container := fmt.Sprintf("eu.gcr.io/altipla-tools/%s:latest", containerName)

	root, err := os.Getwd()
	if err != nil {
		fmt.Println("cannot obtain current working directory:", err.Error())
		os.Exit(3)
	}

	switch {
	case collections.HasString(containers, containerName):
		sh := []string{
			"run", "--rm",
			"-it",
			"--user", fmt.Sprintf("%d:%d", os.Getuid(), os.Getgid()),
			"-v", fmt.Sprintf("%s:/workspace", root),
		}

		sshAuthSock := os.Getenv("SSH_AUTH_SOCK")
		if sshAuthSock == "" {
			fmt.Println("WARNING: No SSH_AUTH_SOCK defined in the environment. Start an ssh-agent to share the SSH keys with the tools.")
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
			fmt.Println("cannot check network:", err.Error())
			os.Exit(3)
		}
		if hasNetwork {
			sh = append(sh, fmt.Sprintf("--network=%s", networkName))
		}

		if hasConfig("glide.yaml") {
			f, err := os.Open("glide.yaml")
			if err != nil {
				fmt.Println("cannot open glide config:", err.Error())
				os.Exit(3)
			}

			content, err := ioutil.ReadAll(f)
			if err != nil {
				fmt.Println("cannot read glide config:", err.Error())
				os.Exit(3)
			}

			cnf := new(glideConfig)
			if err := yaml.Unmarshal(content, cnf); err != nil {
				fmt.Println("cannot decode glide config:", err.Error())
				os.Exit(3)
			}

			if cnf.Package != "." && cnf.Package != "" {
				sh = append(sh, "-v", fmt.Sprintf("%s:/go/src/%s", root, cnf.Package))
				sh = append(sh, "-w", fmt.Sprintf("/go/src/%s", cnf.Package))
			}
		}

		sh = append(sh, container)

		if os.Args[1] == "debug" {
			sh = append(sh, os.Args[3:]...)
		} else {
			sh = append(sh, os.Args[1:]...)
		}

		if err := runCmd("docker", sh...); err != nil {
			fmt.Println("error running tool:", err.Error())
			os.Exit(2)
		}

	case os.Args[1] == "pull":
		if err := pullContainers(); err != nil {
			fmt.Println("error running pull:", err.Error())
			os.Exit(2)
		}

	default:
		fmt.Println("unknown tool")
		os.Exit(3)
	}
}

func runCmd(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}

func pullContainers() error {
	for _, tool := range containers {
		if err := runCmd("docker", "pull", fmt.Sprintf("eu.gcr.io/altipla-tools/%s:latest", tool)); err != nil {
			return err
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
