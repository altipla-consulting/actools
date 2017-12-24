package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Pass the tool name as the first argument: actools glide")
		os.Exit(1)
	}

	container := fmt.Sprintf("eu.gcr.io/altipla-tools/%s:latest", os.Args[1])
	root, err := os.Getwd()
	if err != nil {
		fmt.Println("cannot obtain current working directory:", err.Error())
		os.Exit(3)
	}
	switch os.Args[1] {
	case "glide", "go":
		if err := runGeneric(root, container, os.Args[1:]); err != nil {
			fmt.Println("error running tool:", err.Error())
			os.Exit(2)
		}

	case "pull":
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

func runGeneric(root, container string, args []string) error {
	sh := []string{
		"run", "--rm",
		"-it",
		"--user", fmt.Sprintf("%d:%d", os.Getuid(), os.Getgid()),
		"-v", fmt.Sprintf("%s:/staging", root),
		container,
	}
	sh = append(sh, args...)

	if err := runCmd("docker", sh...); err != nil {
		return err
	}

	return nil
}

func pullContainers() error {
	tools := []string{
		"glide",
		"go",
	}
	for _, tool := range tools {
		if err := runCmd("docker", "pull", fmt.Sprintf("eu.gcr.io/altipla-tools/%s:latest", tool)); err != nil {
			return err
		}
	}
	return nil
}
