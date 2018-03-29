package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"

	"github.com/altipla-consulting/collections"
	"github.com/juju/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/altipla-consulting/actools/pkg/docker"
)

func init() {
	CmdRoot.AddCommand(CmdStart)
}

var CmdStart = &cobra.Command{
	Use:   "start",
	Short: "Enciende un servicio de desarrollo",
	RunE: func(cmd *cobra.Command, args []string) error {
		cnf, err := ReadConfig()
		if err != nil {
			return errors.Trace(err)
		}
		if cnf == nil {
			return errors.NotFoundf("actools.yml")
		}

		root, err := os.Getwd()
		if err != nil {
			return errors.Trace(err)
		}

		networkName := fmt.Sprintf("%s_default", filepath.Base(root))
		hasNetwork, err := docker.NetworkExists(networkName)
		if err != nil {
			return errors.Trace(err)
		}
		if !hasNetwork {
			if err := docker.CreateNetwork(networkName); err != nil {
				return errors.Trace(err)
			}
		}

		start := []string{}
		foreground := []string{}
		for _, arg := range args {
			service, ok := cnf.Services[arg]
			if ok {
				start = append(start, service.Deps...)
				start = append(start, arg)
				foreground = append(foreground, arg)
				continue
			}

			tool, ok := cnf.Tools[arg]
			if ok {
				start = append(start, tool.Deps...)
				start = append(start, arg)
				continue
			}

			return errors.NotFoundf("service %s", arg)
		}

		start = collections.UniqueStrings(start)

		wg := new(sync.WaitGroup)
		notifyExit := make(chan struct{})
		for _, arg := range start {
			name := fmt.Sprintf("%s_%s", filepath.Base(root), arg)
			hasContainer, err := docker.ContainerExists(name)
			if err != nil {
				return errors.Trace(err)
			}

			var container string
			containerArgs := []string{}
			containerCnf := &containerConfig{
				Name:         name,
				Persistent:   true,
				CreateOnly:   true,
				NetworkAlias: []string{arg},
			}

			service, ok := cnf.Services[arg]
			if ok {
				container = fmt.Sprintf("dev-%s", service.Type)

				containerCnf.Ports = service.Ports
				containerCnf.Volumes = service.Volumes
				containerCnf.Env = service.Env
				containerCnf.ShareWorkspace = true
				containerCnf.LocalUser = true
				containerCnf.ShareGcloudConfig = true

				switch service.Type {
				case "go":
					containerCnf.ConfigureGopath = true
					containerArgs = append(containerArgs, "rerun", filepath.Join(cnf.Project, service.Workdir, "cmd", arg))
					containerCnf.Workdir = fmt.Sprintf("/%s", service.Workdir)

				case "gulp":
					containerCnf.Workdir = fmt.Sprintf("/workspace/%s", service.Workdir)
				}
			}

			tool, ok := cnf.Tools[arg]
			if ok {
				container = tool.Container

				containerCnf.Ports = tool.Ports
				containerCnf.Volumes = tool.Volumes
				containerArgs = append(containerArgs, tool.Args)

				switch tool.Container {
				}
			}

			log.WithFields(log.Fields{
				"service": arg,
				"is-new":  !hasContainer,
			}).Info("Start service")

			if !hasContainer {
				if err := runContainer(container, containerCnf, containerArgs...); err != nil {
					return errors.Annotatef(err, "service %s", arg)
				}
			}

			if collections.HasString(foreground, arg) {
				wg.Add(1)
				go RunForeground(wg, notifyExit, arg, name)
			} else {
				if err := runInteractiveDebugOutput("docker", "start", name); err != nil {
					return errors.Trace(err)
				}
			}
		}

		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt)
		go func() {
			for range c {
				// Newline to jump the next log we emit.
				fmt.Println()

				// Enter a loop until all the containers exits and the main closes the whole app directly.
				for {
					notifyExit <- struct{}{}
				}
			}
		}()

		// Wait for all running services to finish.
		wg.Wait()

		return nil
	},
}

type PrefixFormatter struct {
	t           *log.TextFormatter
	serviceName string
}

func (f *PrefixFormatter) Format(entry *log.Entry) ([]byte, error) {
	entry.Message = fmt.Sprintf("(%s) %s", f.serviceName, entry.Message)

	return f.t.Format(entry)
}

func newPrefixFormatter(serviceName string) *PrefixFormatter {
	return &PrefixFormatter{
		t:           new(log.TextFormatter),
		serviceName: serviceName,
	}
}

func RunForeground(wg *sync.WaitGroup, notifyExit chan struct{}, serviceName, containerName string) {
	defer wg.Done()

	logger := log.WithFields(log.Fields{"service": serviceName})

	logger.Debug("Run foreground service")
	notifyErr := make(chan error, 1)
	go func() {
		cmd := exec.Command("docker", "start", "-a", containerName)
		cmd.Stdin = os.Stdin

		loggerOut := log.New()
		loggerOut.Formatter = newPrefixFormatter(serviceName)
		loggerOut.SetLevel(log.StandardLogger().Level)
		wout := &logrusWriter{loggerOut}
		cmd.Stdout = wout

		loggerErr := log.New()
		loggerErr.Formatter = newPrefixFormatter(serviceName)
		loggerErr.SetLevel(log.StandardLogger().Level)
		werr := loggerErr.WriterLevel(log.DebugLevel)
		defer werr.Close()
		cmd.Stderr = werr

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
		logger.WithFields(log.Fields{"err": err.Error()}).Error("Foreground service failed")

	case <-notifyExit:
		logger.Info("Kill service")
		if err := runInteractiveDebugOutput("docker", "kill", containerName); err != nil {
			logger.WithFields(log.Fields{"err": err.Error()}).Error("Stop foreground service failed")
		}
	}
}

type logrusWriter struct {
	logger *log.Logger
}

func (w *logrusWriter) Write(b []byte) (int, error) {
	s := string(b)
	s = strings.TrimSpace(s)
	parts := strings.Split(s, "\r\n")
	for _, part := range parts {
		w.logger.Info(part)
	}

	return len(b), nil
}
