package docker

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/juju/errors"

	"github.com/altipla-consulting/actools/pkg/config"
	"github.com/altipla-consulting/actools/pkg/run"
)

type ContainerManager struct {
	name         string
	image        *ImageManager
	network      *NetworkManager
	localUser    bool
	networkAlias string
	noTTY        bool
	env          map[string]string
	volumes      map[string]string
	ports        map[int64]int64

	// userWorkdir will overwrite workdir if specified
	workdir     string
	userWorkdir string

	persistent bool
}

func Container(name string, options ...ContainerOption) (*ContainerManager, error) {
	container := &ContainerManager{
		name:  fmt.Sprintf("%s_%s", config.ProjectName(), name),
		noTTY: config.Jenkins(),
		env: map[string]string{
			"HOME": "/home/container",
		},
		volumes: make(map[string]string),
		ports:   make(map[int64]int64),
	}

	for _, option := range options {
		if err := option(container); err != nil {
			return nil, errors.Trace(err)
		}
	}

	if container.userWorkdir != "" {
		container.workdir = container.userWorkdir
	}

	return container, nil
}

func (container *ContainerManager) String() string {
	return container.name
}

func (container *ContainerManager) Exists() (bool, error) {
	cmd := exec.Command("docker", "inspect", container.name)
	if err := cmd.Run(); err != nil {
		if !cmd.ProcessState.Success() {
			return false, nil
		}

		return false, errors.Trace(err)
	}

	return true, nil
}

func (container *ContainerManager) Running() (bool, error) {
	exists, err := container.Exists()
	if err != nil {
		return false, errors.Trace(err)
	} else if !exists {
		return false, nil
	}

	cmd := exec.Command("docker", "inspect", "-f", "{{.State.Running}}", container.name)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false, errors.Annotate(err, "cannot inspect container")
	}

	switch strings.TrimSpace(string(output)) {
	case "true":
		return true, nil

	case "false":
		return false, nil
	}

	return false, errors.Errorf("unknown inspect output: %s", output)
}

func (container *ContainerManager) Stop() error {
	if exists, err := container.Exists(); err != nil {
		return errors.Trace(err)
	} else if !exists {
		return nil
	}

	return errors.Trace(run.Interactive("docker", "stop", container.name))
}

func (container *ContainerManager) Start(args ...string) error {
	if exists, err := container.Exists(); err != nil {
		return errors.Trace(err)
	} else if !exists {
		if err := container.Create(args...); err != nil {
			return errors.Trace(err)
		}
	}

	return errors.Trace(run.Interactive("docker", "start", container.name))
}

func (container *ContainerManager) Kill() error {
	exists, err := container.Exists()
	if err != nil {
		return errors.Trace(err)
	}
	if !exists {
		return nil
	}

	return errors.Trace(run.Interactive("docker", "kill", container.name))
}

func (container *ContainerManager) Remove() error {
	exists, err := container.Exists()
	if err != nil {
		return errors.Trace(err)
	}
	if !exists {
		return nil
	}

	return errors.Trace(run.Interactive("docker", "rm", container.name))
}

func (container *ContainerManager) Run(args ...string) error {
	sh, err := container.buildCommand("run", args...)
	if err != nil {
		return errors.Trace(err)
	}

	return errors.Trace(run.InteractiveWithOutput("docker", sh...))
}

func (container *ContainerManager) buildCommand(operation string, args ...string) ([]string, error) {
	// Creamos la red del contenedor si no existía previamente
	if container.network != nil {
		if err := container.network.CreateIfNotExists(); err != nil {
			return nil, errors.Trace(err)
		}
	}

	wd, err := os.Getwd()
	if err != nil {
		return nil, errors.Trace(err)
	}

	var sh []string

	// Ejecutamos interactivamente.
	sh = append(sh, operation, "-i")

	// La terminal solo la podemos activar en local cuando ejecutamos comandos directamente.
	if !container.noTTY {
		sh = append(sh, "-t")
	}

	// Nombre del contenedor.
	sh = append(sh, "--name", container.name)

	// Los contenedores que ejecutamos normalmente son transitorios, excepto los
	// servicios que conservan su estado hasta que son explícitamente reiniciados.
	if !container.persistent {
		sh = append(sh, "--rm")
	}

	// Cambiamos el usuario de dentro para que coincida con el de fuera y los
	// archivos escritos mantengan los permisos iguales en todas partes. En Windows
	// no es necesario puesto que usa Samba y una máquina virtual.
	if container.localUser && config.Linux() {
		sh = append(sh, "--user", fmt.Sprintf("%d:%d", os.Getuid(), os.Getgid()))
	}

	// Configuramos el alias del contenedor en la red local para comunicarnos con él.
	if container.networkAlias != "" {
		sh = append(sh, "--network-alias", container.networkAlias)
	}

	// Variables de entorno del contenedor.
	for k, v := range container.env {
		sh = append(sh, "-e", fmt.Sprintf("%v=%v", k, v))
	}

	// Red en la que se ejecutará el contenedor y que permitirá con el DNS interno
	// comunicarse a los servicios los unos con los otros.
	if container.network != nil {
		sh = append(sh, "--network", container.network.String())
	}

	// Compartimos los puertos con la máquina.
	for source, inside := range container.ports {
		sh = append(sh, "-p", fmt.Sprintf("%d:%d", source, inside))
	}

	for source, inside := range container.volumes {
		// Permitimos acortar las direcciones relativas al directorio actual
		// para facilitar la portabilidad del fichero actools.yml.
		if strings.HasPrefix(source, "./") {
			source = filepath.Join(wd, source[2:])
		}

		// Compartimos los volúmenes del contenedor.
		sh = append(sh, "-v", fmt.Sprintf("%s:%s", source, inside))
	}

	// Directorio de trabajo inicial al abrir el contenedor.
	if container.workdir != "" {
		sh = append(sh, "-w", container.workdir)
	}

	// Añadimos la imagen que ejecutamos.
	sh = append(sh, container.image.String())

	// Añadimos cualquier adicional que recibamos en el momento.
	sh = append(sh, args...)

	return sh, nil
}

func (container *ContainerManager) Create(args ...string) error {
	if exists, err := container.Exists(); err != nil {
		return errors.Trace(err)
	} else if exists {
		return nil
	}

	sh, err := container.buildCommand("create", args...)
	if err != nil {
		return errors.Trace(err)
	}

	return errors.Trace(run.Interactive("docker", sh...))
}
