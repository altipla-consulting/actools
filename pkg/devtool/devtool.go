package devtool

import (
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"libs.altipla.consulting/errors"

	"github.com/altipla-consulting/actools/pkg/config"
	"github.com/altipla-consulting/actools/pkg/run"
)

var (
	compileCh = make(chan string)
	runtimeCh = make(chan string)
)

func Run(apps []string) error {
	var g errgroup.Group
	g.Go(compiler)
	g.Go(sourceCodeWatcher(apps))

	for _, app := range apps {
		compileCh <- app
	}

	return errors.Trace(g.Wait())
}

func compiler() error {
	pending := map[string]bool{}
	var timer *time.Timer
	var timerCh <-chan time.Time

	for {
		timerCh = nil
		if timer != nil {
			timerCh = timer.C
		}

		select {
		case app := <-compileCh:
			pending[app] = true

			if timer != nil && !timer.Stop() {
				<-timer.C
			}
			timer = time.NewTimer(3000 * time.Millisecond)

		case <-timerCh:
			timer = nil
			for app := range pending {
				delete(pending, app)

				log.WithField("app", app).Info("Building app")

				svc := config.Settings.Services[app]
				if err := run.InteractiveWithOutput("go", "install", "./"+svc.Workdir); err != nil {
					exit, ok := errors.Cause(err).(*exec.ExitError)
					if !ok {
						return errors.Trace(err)
					}

					log.WithFields(log.Fields{
						"app":       app,
						"exit-code": exit.ExitCode(),
					}).Error("App compilation failed")
				}
			}
		}
	}

	return nil
}

func sourceCodeWatcher(apps []string) func() error {
	mapping := map[string]string{}

	return func() error {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			return errors.Trace(err)
		}

		for _, app := range apps {
			svc := config.Settings.Services[app]
			walkFn := func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return errors.Trace(err)
				}
				if !info.IsDir() {
					return nil
				}

				log.WithField("path", path).Debug("Watching directory for changes")
				mapping[path] = app
				if err := watcher.Add(path); err != nil {
					return errors.Trace(err)
				}

				return nil
			}
			if err := filepath.Walk(svc.Workdir, walkFn); err != nil {
				return errors.Trace(err)
			}
		}

		for {
			select {
			case err := <-watcher.Errors:
				return errors.Trace(err)

			case ev := <-watcher.Events:
				switch ev.Op {
				case fsnotify.Create:
					app := mapping[filepath.Dir(ev.Name)]

					info, err := os.Stat(ev.Name)
					if err != nil {
						return errors.Trace(err)
					}
					if info.IsDir() {
						log.WithField("path", ev.Name).Info("New directory detected and monitored for changes")
						mapping[ev.Name] = app
						if err := watcher.Add(ev.Name); err != nil {
							return errors.Trace(err)
						}
					} else {
						log.WithFields(log.Fields{
							"path": ev.Name,
							"app":  app,
						}).Debug("New file detected")
						if err := processChange(app, ev.Name); err != nil {
							return errors.Trace(err)
						}
					}

				case fsnotify.Write:
					app := mapping[filepath.Dir(ev.Name)]

					log.WithFields(log.Fields{
						"path": ev.Name,
						"app":  app,
					}).Debug("File write detected")
					if err := processChange(app, ev.Name); err != nil {
						return errors.Trace(err)
					}

				case fsnotify.Remove:
					if mapping[ev.Name] != "" {
						log.WithFields(log.Fields{
							"path": ev.Name,
							"app":  mapping[ev.Name],
						}).Debug("Removed directory detected")
						delete(mapping, ev.Name)
					} else {
						app := mapping[filepath.Dir(ev.Name)]

						log.WithFields(log.Fields{
							"path": ev.Name,
							"app":  app,
						}).Debug("Removed file detected")

						if err := processChange(app, ev.Name); err != nil {
							return errors.Trace(err)
						}
					}

					// TODO(ernesto): Debemos recoger renames?
				}
			}
		}

		return nil
	}
}

func processChange(app, name string) error {
	if filepath.Ext(name) == ".go" {
		compileCh <- app
	}
	
	return nil
}
