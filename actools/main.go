package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/altipla-consulting/collections"
)

var containers = []string{
	"glide",
	"go",
	"node",
	"protoc",
	"rambler",
}

var toolContainers = map[string]string{
	"npm": "node",
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
	container := fmt.Sprintf("eu.gcr.io/altipla-tools/%s:latest", containerName)

	root, err := os.Getwd()
	if err != nil {
		fmt.Println("cannot obtain current working directory:", err.Error())
		os.Exit(3)
	}

	switch {
	case collections.HasString(containers, containerName):
		sshAuthSock := os.Getenv("SSH_AUTH_SOCK")
		if sshAuthSock == "" {
			fmt.Println("WARNING: No SSH_AUTH_SOCK defined in the environment. Start an ssh-agent to share the SSH keys with the tools.")
		}

		hasGcloudConfig := true
		gcloudConfigPath := fmt.Sprintf("%s/.config/gcloud", os.Getenv("HOME"))
		if _, err := os.Stat(gcloudConfigPath); err != nil {
			if !os.IsNotExist(err) {
				fmt.Println("cannot stat the gcloud config path:", err.Error)
				os.Exit(3)
			}

			hasGcloudConfig = false
		}

		sh := []string{
			"run", "--rm",
			"-it",
			"--user", fmt.Sprintf("%d:%d", os.Getuid(), os.Getgid()),
			"-v", fmt.Sprintf("%s:/staging", root),
		}
		if sshAuthSock != "" {
			sh = append(sh, "-e", fmt.Sprintf("SSH_AUTH_SOCK=%s", sshAuthSock))
			sh = append(sh, "-v", fmt.Sprintf("%s:%s", sshAuthSock, sshAuthSock))
		}
		if hasGcloudConfig {
			sh = append(sh, "-v", fmt.Sprintf("%s:/root/.config/gcloud", gcloudConfigPath))
		}
		sh = append(sh, container)
		sh = append(sh, os.Args[1:]...)

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
