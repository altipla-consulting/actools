package run

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"io"
	"bufio"

	log "github.com/sirupsen/logrus"
	"libs.altipla.consulting/errors"
)

func logCommand(name string, args []string) {
	// Do not use fields as they escape the value and the result cannot be copied.
	logArgs := []string{}
	for _, arg := range args {
		if strings.Contains(arg, `"`) {
			logArgs = append(logArgs, strconv.Quote(arg))
		} else {
			logArgs = append(logArgs, arg)
		}
	}
	log.Debugln(fmt.Sprintf("%s %s", name, strings.Join(logArgs, " ")))
}

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

	log.Debug("Run interactive command")
	logCommand(name, args)

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

	log.Debug("Run interactive command with output")
	logCommand(name, args)

	if err := cmd.Run(); err != nil {
		return errors.Trace(err)
	}

	return nil
}

func NonInteractiveWithOutput(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Debug("Run non interactive command with output")
	logCommand(name, args)

	if err := cmd.Run(); err != nil {
		return errors.Trace(err)
	}

	return nil
}

func NonInteractiveCaptureOutput(linesToCapture int, name string, args ...string) ([]string, error) {
	reader, writer := io.Pipe()
	defer reader.Close()
	defer writer.Close()

	lines := make([]string, 0, linesToCapture)

	go func() {
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			line := scanner.Text()

			if len(lines) < linesToCapture {
				lines = append(lines, line)
			} else {
				for i := 0; i < linesToCapture - 1; i++ {
					lines[i] = lines[i+1]
				}
				lines[linesToCapture-1] = line
			}

			fmt.Fprintln(os.Stderr, line)
		}
		if err := scanner.Err(); err != nil && !errors.Is(err, io.ErrClosedPipe) {
			panic(err)
		}
	}()

	cmd := exec.Command(name, args...)
	cmd.Stdout = writer
	cmd.Stderr = writer

	log.Debug("Run non interactive command capturing output")
	logCommand(name, args)

	if err := cmd.Run(); err != nil {
		return lines, errors.Trace(err)
	}

	return lines, nil
}

func InteractiveCaptureOutput(name string, args ...string) (string, error) {
	var buf bytes.Buffer

	cmd := exec.Command(name, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = &buf
	cmd.Stderr = os.Stderr

	log.Debug("Run interactive command capturing output")
	logCommand(name, args)

	if err := cmd.Run(); err != nil {
		return "", errors.Trace(err)
	}

	return strings.TrimSpace(buf.String()), nil
}
