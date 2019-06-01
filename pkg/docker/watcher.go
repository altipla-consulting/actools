package docker

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"sync"

	log "github.com/sirupsen/logrus"
)

type Watcher struct {
	wg         *sync.WaitGroup
	notifyExit chan struct{}
}

func NewWatcher() *Watcher {
	return &Watcher{
		wg:         new(sync.WaitGroup),
		notifyExit: make(chan struct{}),
	}
}

func (watcher *Watcher) Run(serviceName string, container *ContainerManager) {
	watcher.wg.Add(1)
	go runForeground(watcher.wg, watcher.notifyExit, serviceName, container)
}

func (watcher *Watcher) Wait() {
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt)
		for range c {
			// Newline to always jump the next log we emit.
			fmt.Println()

			// Enter a loop until all the containers exits and the main closes the whole app directly.
			for {
				watcher.notifyExit <- struct{}{}
			}
		}
	}()

	// Wait for all running services to finish.
	watcher.wg.Wait()
}

func runForeground(wg *sync.WaitGroup, notifyExit chan struct{}, serviceName string, container *ContainerManager) {
	defer wg.Done()

	logger := log.WithField("service", serviceName)
	logger.Info("Start service")

	notifyErr := make(chan error, 1)
	go func() {
		if err := container.Create(); err != nil {
			notifyErr <- err
			return
		}

		cmd := exec.Command("docker", "start", "-a", container.String())
		cmd.Stdin = os.Stdin

		loggerOut := log.New()
		loggerOut.Formatter = newPrefixFormatter(serviceName)
		loggerOut.SetLevel(log.StandardLogger().Level)
		cmd.Stdout = &logrusWriter{logger: loggerOut}

		loggerErr := log.New()
		loggerErr.Formatter = newPrefixFormatter(serviceName)
		loggerErr.SetLevel(log.StandardLogger().Level)
		cmd.Stderr = &logrusWriter{logger: loggerErr, debug: true}

		if err := cmd.Start(); err != nil {
			notifyErr <- err
			return
		}
		if err := cmd.Wait(); err != nil {
			if err.Error() == "signal: interrupt" {
				return
			}

			notifyErr <- err
			return
		}
	}()

	select {
	case err := <-notifyErr:
		logger.WithField("err", err.Error()).Error("Service failed")

	case <-notifyExit:
		logger.Info("Kill service")
		if err := container.Kill(); err != nil {
			logger.WithField("err", err.Error()).Error("Stop service failed")
		}
	}
}

type prefixFormatter struct {
	t           *log.TextFormatter
	serviceName string
}

func (f *prefixFormatter) Format(entry *log.Entry) ([]byte, error) {
	entry.Message = fmt.Sprintf("(%s) %s", f.serviceName, entry.Message)

	return f.t.Format(entry)
}

func newPrefixFormatter(serviceName string) *prefixFormatter {
	return &prefixFormatter{
		t:           new(log.TextFormatter),
		serviceName: serviceName,
	}
}

type logrusWriter struct {
	logger *log.Logger
	debug  bool
}

func (w *logrusWriter) Write(b []byte) (int, error) {
	const maxLogSize = 64 * 1024

	var result [][]byte
	parts := bytes.Split(bytes.TrimSpace(b), []byte("\r\n"))
	for _, part := range parts {
		for len(part) > maxLogSize {
			result = append(result, part[:maxLogSize])
			part = part[maxLogSize:]
		}

		if len(part) > 0 {
			result = append(result, part)
		}
	}
	for _, line := range result {
		if w.debug {
			w.logger.Debug(string(line))
		} else {
			w.logger.Info(string(line))
		}
	}

	return len(b), nil
}
