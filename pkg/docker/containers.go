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
	workdir      string
}

func Container(name string, options ...ContainerOption) (*ContainerManager, error) {
	container := &ContainerManager{
		name:  fmt.Sprintf("%s_%s", config.ProjectName(), name),
		noTTY: config.Jenkins(),
		env:   make(map[string]string),
		ports: make(map[int64]int64),
	}

	for _, option := range options {
		if err := option(container); err != nil {
			return nil, errors.Trace(err)
		}
	}

	return container, nil
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

func (container *ContainerManager) Stop() error {
	exists, err := container.Exists()
	if err != nil {
		return errors.Trace(err)
	}
	if !exists {
		return nil
	}

	return errors.Trace(run.Interactive("docker", "stop", container.name))
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

func (container *ContainerManager) Run() error {
	// Creamos la red del contenedor si no existía previamente
	if container.network != nil {
		if err := container.network.CreateIfNotExists(); err != nil {
			return errors.Trace(err)
		}
	}

	wd, err := os.Getwd()
	if err != nil {
		return errors.Trace(err)
	}

	var args []string

	// Ejecutamos interactivamente.
	args = append(args, "run", "-i")

	// La terminal solo la podemos activar en local cuando ejecutamos comandos directamente.
	if !container.noTTY {
		args = append(args, "-t")
	}

	// Nombre del contenedor.
	args = append(args, "--name", container.name)

	// Los contenedores que ejecutamos siempre son transitorios.
	args = append(args, "--rm")

	// Cambiamos el usuario de dentro para que coincida con el de fuera y los
	// archivos escritos mantengan los permisos iguales en todas partes. En Windows
	// no es necesario puesto que usa Samba y una máquina virtual.
	if container.localUser && config.Linux() {
		args = append(args, "--user", fmt.Sprintf("%d:%d", os.Getuid(), os.Getgid()))
	}

	// Configuramos el alias del contenedor en la red local para comunicarnos con él.
	if container.networkAlias != "" {
		args = append(args, "--network-alias", container.networkAlias)
	}

	// Variables de entorno del contenedor.
	for k, v := range container.env {
		args = append(args, "-e", fmt.Sprintf("%v=%v", k, v))
	}

	// Red en la que se ejecutará el contenedor y que permitirá con el DNS interno
	// comunicarse a los servicios los unos con los otros.
	if container.network != nil {
		args = append(args, "--network", container.network.String())
	}

	// Compartimos los puertos con la máquina.
	for source, inside := range container.ports {
		args = append(args, "-p", fmt.Sprintf("%d:%d", source, inside))
	}

	for source, inside := range container.volumes {
		// Permitimos acortar las direcciones relativas al directorio actual
		// para facilitar la portabilidad del fichero actools.yml.
		if strings.HasPrefix(source, "./") {
			source = filepath.Join(wd, source[2:])
		}

		// Compartimos los volúmenes del contenedor.
		args = append(args, "-v", fmt.Sprintf("%s:%s", source, inside))
	}

	// Directorio de trabajo inicial al abrir el contenedor.
	if container.workdir != "" {
		args = append(args, "-w", container.workdir)
	}

	// Añadimos la imagen que ejecutamos.
	args = append(args, container.image.String())

	return errors.Trace(run.Interactive("docker", args...))
}
