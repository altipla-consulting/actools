package config

import (
	"os"
	"runtime"
)

func Jenkins() bool {
	return os.Getenv("JENKINS_URL") != ""
}

func SSHAgentSocket() string {
	return os.Getenv("SSH_AUTH_SOCK")
}

func Linux() bool {
	return runtime.GOOS == "linux"
}

func Home() string {
	return os.Getenv("HOME")
}

func Development() bool {
	return Version == "dev"
}
