package devtool

import (
	"context"
	"fmt"
	"os"
	"os/exec"
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

const changeDelay = 1 * time.Second

var (
	compileCh = make(chan string)
)

func Run(ctx context.Context, g *errgroup.Group, watcher *docker.Watcher, apps []string) error {
	restartChs := make(map[string]chan struct{})
	for _, app := range apps {
		restartCh := make(chan struct{})
		restartChs[app] = restartCh

		g.Go(sourceCodeWatcher(ctx, app))
		g.Go(configWatcher(ctx, app, restartCh))
		g.Go(serviceRunner(ctx, app, watcher, restartCh))
	}

	g.Go(compiler(ctx, restartChs))

	// Initial compilation of all the apps sequentially.
	for _, app := range apps {
		compileCh <- app
	}

	return nil
}

type compilationFailure struct {
	app    string
	output []string
}

func compiler(ctx context.Context, restartChs map[string]chan struct{}) func() error {
	return func() error {
		pending := map[string]bool{}
		var timer *time.Timer
		var timerCh <-chan time.Time
		wasFailing := map[string]bool{}
		working := make(map[string]bool)

		sendCompilationCh := make(chan string)
		successCh := make(chan string)
		failureCh := make(chan *compilationFailure)
		for i := 0; i < 6; i++ {
			go compilerWorker(ctx, i, sendCompilationCh, successCh, failureCh)
		}

		first := true
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
				timer = time.NewTimer(changeDelay)

			case <-timerCh:
				timer = nil

				// El flag global de failed hace que pasemos a considerar el trabajo
				// ya lanzado y terminemos las compilaciones sin mandar ninguna más.
				var failed bool
				for (!failed && len(pending) > 0) || len(working) > 0 {
					// Cogemos la primera compilación pendiente que encontremos, siempre
					// que no haya fallado ya el proceso globalmente.
					var first string
					var ch chan string
					if !failed {
						for app := range pending {
							ch = sendCompilationCh
							first = app
							break
						}
					}

					select {
					// Mandamos una nueva compilación a un worker satisfactoriamente.
					case ch <- first:
						delete(pending, first)
						working[first] = true

					// El worker ha compilado correctamente su trabajo.
					case app := <-successCh:
						delete(working, app)

						if wasFailing[app] {
							wasFailing[app] = false
							if err := notify.Send(app, "Compilado correctamente.", notify.IconInfo); err != nil {
								return errors.Trace(err)
							}
						}

						restartChs[app] <- struct{}{}

					// El worker no ha completar la compilación sin errores.
					case fail := <-failureCh:
						failed = true
						delete(working, fail.app)
						wasFailing[fail.app] = true

						log.WithField("app", fail.app).Error("App compilation failed")

						if len(fail.output) > 0 {
							if err := notify.Send(fail.app, strings.Join(fail.output, "\n"), notify.IconError); err != nil {
								return errors.Trace(err)
							}
						}
					}
				}

				if first {
					first = false
					log.Warning("")
					log.Warning("----------------------------------------------")
					log.Warning("--- ALL COMPILATIONS FINISHED SUCCESSFULLY ---")
					log.Warning("----------------------------------------------")
					log.Warning("")
				}
			}
		}
	}

	return nil
}

func compilerWorker(ctx context.Context, worker int, sendCompilationCh <-chan string, successCh chan<- string, failureCh chan<- *compilationFailure) {
	tool := containers.FindContainerTool("go")

	options := []docker.ContainerOption{
		docker.WithImage(docker.Image(containers.Repo, tool.Image)),
		docker.WithDefaultNetwork(),
	}
	options = append(options, tool.Options...)

	container, err := docker.Container(fmt.Sprintf("devtool-%s-go-%d", tool.Image, worker), options...)
	if err != nil {
		panic(err)
	}

	for {
		select {
		case <-ctx.Done():
			return

		case app := <-sendCompilationCh:
			log.WithField("app", app).Info("Building app")

			svc := config.Settings.Services[app]
			lines, err := container.RunNonInteractiveCaptureOutput(7, "go", "install", "./"+svc.Workdir)
			if err != nil {
				if _, ok := errors.Cause(err).(*exec.ExitError); ok {
					failureCh <- &compilationFailure{
						app:    app,
						output: lines,
					}
				} else {
					failureCh <- &compilationFailure{
						app:    app,
						output: append(lines, "INTERNAL ERROR: "+err.Error()),
					}
				}
			} else {
				successCh <- app
			}
		}
	}
}

func sourceCodeWatcher(ctx context.Context, app string) func() error {
	return func() error {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			return errors.Trace(err)
		}

		svc := config.Settings.Services[app]
		walkFn := func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return errors.Trace(err)
			}
			if !info.IsDir() {
				return nil
			}

			var ignore bool
			for _, ig := range svc.Ignore {
				if strings.HasPrefix(path, ig) {
					ignore = true
					break
				}
			}
			if ignore {
				return nil
			}

			log.WithField("path", path).Debug("Watching directory for changes")
			if err := watcher.Add(path); err != nil {
				return errors.Trace(err)
			}

			return nil
		}

		// Watch app source code folder.
		if err := filepath.Walk(svc.Workdir, walkFn); err != nil {
			return errors.Trace(err)
		}

		// Watch a root pkg folder with Go libs if it exists.
		if _, err := os.Stat("pkg"); err != nil && !os.IsNotExist(err) {
			return errors.Trace(err)
		} else if err == nil {
			if err := filepath.Walk("pkg", walkFn); err != nil {
				return errors.Trace(err)
			}
		}

		for {
			select {
			case <-ctx.Done():
				return nil

			case err := <-watcher.Errors:
				return errors.Trace(err)

			case ev := <-watcher.Events:
				var ignore bool
				for _, ig := range svc.Ignore {
					if strings.HasPrefix(ev.Name, ig) {
						ignore = true
						break
					}
				}
				if ignore {
					continue
				}

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
						}).Debug("New source code file detected")
						if err := processChange(app, ev.Name); err != nil {
							return errors.Trace(err)
						}
					}

				case fsnotify.Write, fsnotify.Remove, fsnotify.Rename:
					log.WithFields(log.Fields{
						"path": ev.Name,
						"app":  app,
					}).Debug("Source code file change detected")
					if err := processChange(app, ev.Name); err != nil {
						return errors.Trace(err)
					}
				}
			}
		}

		return nil
	}
}

func configWatcher(ctx context.Context, app string, restartCh chan struct{}) func() error {
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
		for _, volume := range svc.Volumes {
			path := strings.Split(volume, ":")[0]
			if err := filepath.Walk(path, walkFn); err != nil {
				return errors.Trace(err)
			}
		}

		var timer *time.Timer
		var timerCh <-chan time.Time
		for {
			timerCh = nil
			if timer != nil {
				timerCh = timer.C
			}

			select {
			case <-ctx.Done():
				return nil

			case err := <-watcher.Errors:
				return errors.Trace(err)

			case ev := <-watcher.Events:
				log.WithFields(log.Fields{
					"path": ev.Name,
					"app":  app,
				}).Debug("Config file change detected")
				if timer != nil && !timer.Stop() {
					<-timer.C
				}
				timer = time.NewTimer(changeDelay)

			case <-timerCh:
				timer = nil
				restartCh <- struct{}{}
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
