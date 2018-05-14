package run

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"
)

func Interactive(name string, args ...string) error {
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

	logArgs := []string{}
	for _, arg := range args {
		logArgs = append(logArgs, strconv.Quote(arg))
	}
	shell := fmt.Sprintf("%s %s", name, strings.Join(logArgs, " "))
	log.WithField("cmd", shell).Debug("Run interactive command")

	if err := cmd.Run(); err != nil {
		return errors.Trace(err)
	}

	return nil
}

func InteractiveWithOutput(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	logArgs := []string{}
	for _, arg := range args {
		logArgs = append(logArgs, strconv.Quote(arg))
	}
	shell := fmt.Sprintf("%s %s", name, strings.Join(logArgs, " "))
	log.WithField("cmd", shell).Debug("Run interactive command with output")

	if err := cmd.Run(); err != nil {
		return errors.Trace(err)
	}

	return nil
}
