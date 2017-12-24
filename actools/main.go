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

	switch  {
	case collections.HasString(containers, containerName):
		sh := []string{
			"run", "--rm",
			"-it",
			"--user", fmt.Sprintf("%d:%d", os.Getuid(), os.Getgid()),
			"-v", fmt.Sprintf("%s:/staging", root),
			container,
		}
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
