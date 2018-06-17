package run

import (
	"bytes"
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

	// Print the command. Do not use fields as they escape the value and the
	// result cannot be copied.
	logArgs := []string{}
	for _, arg := range args {
		if strings.Contains(arg, `"`) {
			logArgs = append(logArgs, strconv.Quote(arg))
		} else {
			logArgs = append(logArgs, arg)
		}
	}
	shell := fmt.Sprintf("%s %s", name, strings.Join(logArgs, " "))
	log.Debug("Run interactive command")
	log.Debugln(shell)

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

	// Print the command. Do not use fields as they escape the value and the
	// result cannot be copied.
	logArgs := []string{}
	for _, arg := range args {
		if strings.Contains(arg, `"`) {
			logArgs = append(logArgs, strconv.Quote(arg))
		} else {
			logArgs = append(logArgs, arg)
		}
	}
	shell := fmt.Sprintf("%s %s", name, strings.Join(logArgs, " "))
	log.Debug("Run interactive command with output")
	log.Debugln(shell)

	if err := cmd.Run(); err != nil {
		return errors.Trace(err)
	}

	return nil
}

func InteractiveCaptureOutput(name string, args ...string) (string, error) {
	var buf bytes.Buffer

	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = &buf
	cmd.Stderr = os.Stderr

	// Print the command. Do not use fields as they escape the value and the
	// result cannot be copied.
	logArgs := []string{}
	for _, arg := range args {
		if strings.Contains(arg, `"`) {
			logArgs = append(logArgs, strconv.Quote(arg))
		} else {
			logArgs = append(logArgs, arg)
		}
	}
	shell := fmt.Sprintf("%s %s", name, strings.Join(logArgs, " "))
	log.Debug("Run interactive command with output")
	log.Debugln(shell)

	if err := cmd.Run(); err != nil {
		return "", errors.Trace(err)
	}

	return strings.TrimSpace(buf.String()), nil
}
