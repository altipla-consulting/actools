package docker

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"libs.altipla.consulting/errors"
)

type Watcher struct {
	sync.Mutex
	g     *errgroup.Group
	stop  map[string]chan struct{}
	ended chan struct{}
}

func NewWatcher() *Watcher {
	return &Watcher{
		g:     new(errgroup.Group),
		stop:  make(map[string]chan struct{}),
		ended: make(chan struct{}),
	}
}

func (watcher *Watcher) Run(serviceName string, container *ContainerManager) {
	watcher.Lock()
	defer watcher.Unlock()

	if watcher.stop[serviceName] != nil {
		panic("repeated service: " + serviceName)
	}
	watcher.stop[serviceName] = make(chan struct{})

	watcher.g.Go(runForeground(watcher.g, watcher.ended, watcher.stop[serviceName], serviceName, container))
}

func (watcher *Watcher) StopAll() error {
	watcher.Lock()
	defer watcher.Unlock()

	for _, ch := range watcher.stop {
		ch <- struct{}{}
		<-watcher.ended
	}

	return errors.Trace(watcher.g.Wait())
}

func (watcher *Watcher) WaitInterrupt() error {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	for range c {
		// Newline to always jump the next log we emit.
		fmt.Println()

		// Enter a loop until all the containers exit.
		for _, ch := range watcher.stop {
			ch <- struct{}{}
		}

		return errors.Trace(watcher.g.Wait())
	}

	return nil
}

func (watcher *Watcher) Stop(serviceName string) {
	watcher.Lock()
	defer watcher.Unlock()

	stopCh, ok := watcher.stop[serviceName]
	if ok {
		stopCh <- struct{}{}
		<-watcher.ended
		delete(watcher.stop, serviceName)
	}
}

func runForeground(g *errgroup.Group, ended, stopCh chan struct{}, serviceName string, container *ContainerManager) func() error {
	return func() error {
		defer func() {
			ended <- struct{}{}
		}()

		reader, writer := io.Pipe()
		defer reader.Close()
		defer writer.Close()

		g.Go(func() error {
			scanner := bufio.NewScanner(reader)
			for scanner.Scan() {
				fmt.Println("(" + serviceName + ") " + scanner.Text())
			}
			if err := scanner.Err(); err != nil && !errors.Is(err, io.ErrClosedPipe) {
				return errors.Trace(err)
			}
			return nil
		})

		logger := log.WithField("service", serviceName)
		logger.Info("Start service")

		if err := container.Create(); err != nil {
			return errors.Trace(err)
		}

		log.Debug("Run container")
		args := []string{"docker", "start", "-a", container.String()}
		log.Debugln(strings.Join(args, " "))

		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = writer
		cmd.Stderr = writer
		if err := cmd.Start(); err != nil {
			return errors.Trace(err)
		}

		<-stopCh

		logger.Info("Stopping service")

		timer := time.NewTimer(5 * time.Second)
		stopped := make(chan bool, 1)
		go func() {
			if err := container.Stop(); err != nil {
				log.WithFields(errors.LogFields(err)).Warning("Cannot stop container")
				return
			}
			cmd.Wait()
			stopped <- true
		}()
		select {
		case <-timer.C:
			if err := container.Kill(); err != nil {
				return errors.Trace(err)
			}
			logger.Warning("Service didn't exited on time and was killed")
		case <-stopped:
			if !timer.Stop() {
				<-timer.C
			}
		}

		return nil
	}
}
