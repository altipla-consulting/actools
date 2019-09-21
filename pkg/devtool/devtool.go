package devtool

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"libs.altipla.consulting/errors"

	"github.com/altipla-consulting/actools/pkg/config"
	"github.com/altipla-consulting/actools/pkg/containers"
	"github.com/altipla-consulting/actools/pkg/docker"
	"github.com/altipla-consulting/actools/pkg/notify"
)

var (
	compileCh = make(chan string)
)

func Run(apps []string) error {
	for _, app := range apps {
		if config.Settings.Services[app] == nil {
			return errors.Errorf("service not found: %s", app)
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	g, ctx := errgroup.WithContext(ctx)

	watcher := docker.NewWatcher()
	g.Go(func() error {
		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt)
		for range c {
			// Newline to always jump the next log we emit.
			fmt.Println()

			if err := watcher.StopAll(); err != nil {
				return errors.Trace(err)
			}

			cancel()
			return nil
		}

		return nil
	})

	restartChs := make(map[string]chan struct{})
	for _, app := range apps {
		restartCh := make(chan struct{})
		restartChs[app] = restartCh

		g.Go(sourceCodeWatcher(ctx, app, restartCh))
		g.Go(serviceRunner(ctx, app, watcher, restartCh))
	}

	g.Go(compiler(ctx, restartChs))

	// Initial compilation of all the apps sequentially.
	for _, app := range apps {
		compileCh <- app
	}

	return errors.Trace(g.Wait())
}

func compiler(ctx context.Context, restartChs map[string]chan struct{}) func() error {
	return func() error {
		pending := map[string]bool{}
		var timer *time.Timer
		var timerCh <-chan time.Time
		wasFailing := map[string]bool{}

		for {
			timerCh = nil
			if timer != nil {
				timerCh = timer.C
			}

			select {
			case <-ctx.Done():
				return nil

			case app := <-compileCh:
				pending[app] = true

				if timer != nil && !timer.Stop() {
					<-timer.C
				}
				timer = time.NewTimer(500 * time.Millisecond)

			case <-timerCh:
				timer = nil
				for app := range pending {
					delete(pending, app)

					log.WithField("app", app).Info("Building app")

					tool := containers.FindContainerTool("go")

					options := []docker.ContainerOption{
						docker.WithImage(docker.Image(containers.Repo, tool.Image)),
						docker.WithDefaultNetwork(),
					}
					options = append(options, tool.Options...)

					container, err := docker.Container(fmt.Sprintf("tool-%s-go", tool.Image), options...)
					if err != nil {
						return errors.Trace(err)
					}

					svc := config.Settings.Services[app]
					if err := container.RunNonInteractive("go", "install", "./"+svc.Workdir); err != nil {
						exit, ok := errors.Cause(err).(*exec.ExitError)
						if !ok {
							return errors.Trace(err)
						}

						log.WithFields(log.Fields{
							"app":       app,
							"exit-code": exit.ExitCode(),
						}).Error("App compilation failed")

						if err := notify.Send(app, "La compilación ha fallado. Mira la terminal para más información.", notify.IconError); err != nil {
							return errors.Trace(err)
						}
						wasFailing[app] = true

						continue
					}

					if wasFailing[app] {
						wasFailing[app] = false
						if err := notify.Send(app, "La compilación va correctamente.", notify.IconInfo); err != nil {
							return errors.Trace(err)
						}
					}

					restartChs[app] <- struct{}{}
				}
			}
		}
	}

	return nil
}

func sourceCodeWatcher(ctx context.Context, app string, restartCh chan struct{}) func() error {
	return func() error {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			return errors.Trace(err)
		}

		walkFn := func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return errors.Trace(err)
			}
			if !info.IsDir() {
				return nil
			}

			log.WithField("path", path).Debug("Watching directory for changes")
			if err := watcher.Add(path); err != nil {
				return errors.Trace(err)
			}

			return nil
		}
		svc := config.Settings.Services[app]
		if err := filepath.Walk(svc.Workdir, walkFn); err != nil {
			return errors.Trace(err)
		}

		for {
			select {
			case <-ctx.Done():
				return nil

			case err := <-watcher.Errors:
				return errors.Trace(err)

			case ev := <-watcher.Events:
				switch ev.Op {
				case fsnotify.Create:
					info, err := os.Stat(ev.Name)
					if err != nil {
						return errors.Trace(err)
					}
					if info.IsDir() {
						log.WithField("path", ev.Name).Info("New directory detected and monitored for changes")
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

				case fsnotify.Write, fsnotify.Remove, fsnotify.Rename:
					log.WithFields(log.Fields{
						"path": ev.Name,
						"app":  app,
					}).Debug("File change detected")
					if err := processChange(app, ev.Name); err != nil {
						return errors.Trace(err)
					}
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

func serviceRunner(ctx context.Context, app string, watcher *docker.Watcher, restartCh chan struct{}) func() error {
	return func() error {
		for {
			select {
			case <-ctx.Done():
				return nil

			case <-restartCh:
				watcher.Stop(app)

				svc := config.Settings.Services[app]
				container, err := containers.FindImage("dev-go")
				if err != nil {
					return errors.Trace(err)
				}

				options := []docker.ContainerOption{
					docker.WithImage(docker.Image(containers.Repo, container.Image)),
					docker.WithDefaultNetwork(),
					docker.WithPersistence(),
					docker.WithWorkdir(fmt.Sprintf("/workspace/%s", svc.Workdir)),
					docker.WithNetworkAlias(app),
					docker.WithEnv("APP", app),
				}
				options = append(options, container.Options...)

				for _, port := range svc.Ports {
					parts := strings.Split(port, ":")
					if len(parts) != 2 {
						return errors.Errorf("invalid ports of service: %s", app)
					}

					source, err := strconv.ParseInt(parts[0], 10, 64)
					if err != nil {
						return errors.Wrapf(err, "invalid port number: %s", parts[0])
					}

					inside, err := strconv.ParseInt(parts[1], 10, 64)
					if err != nil {
						return errors.Wrapf(err, "invalid port number: %s", parts[1])
					}

					options = append(options, docker.WithPort(source, inside))
				}

				for _, volume := range svc.Volumes {
					parts := strings.Split(volume, ":")
					if len(parts) != 2 {
						return errors.Errorf("invalid volumes of service: %s", app)
					}

					options = append(options, docker.WithVolume(parts[0], parts[1]))
				}

				for k, v := range svc.Env {
					options = append(options, docker.WithEnv(k, v))
				}

				cnt, err := docker.Container(app, options...)
				if err != nil {
					return errors.Trace(err)
				}
				watcher.Run(app, cnt)
			}
		}

		return nil
	}
}
