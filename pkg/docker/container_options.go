package docker

import (
	"fmt"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"libs.altipla.consulting/errors"

	"github.com/altipla-consulting/actools/pkg/config"
)

type ContainerOption func(container *ContainerManager) error

func WithNetwork(network *NetworkManager) ContainerOption {
	return func(container *ContainerManager) error {
		container.network = network
		return nil
	}
}

func WithDefaultNetwork() ContainerOption {
	return func(container *ContainerManager) error {
		wd, err := os.Getwd()
		if err != nil {
			return errors.Trace(err)
		}

		container.network = Network(filepath.Base(wd))

		return nil
	}
}

func WithImage(image *ImageManager) ContainerOption {
	return func(container *ContainerManager) error {
		container.image = image
		return nil
	}
}

func WithLocalUser() ContainerOption {
	return func(container *ContainerManager) error {
		container.localUser = true
		return nil
	}
}

func WithSharedSSHSocket() ContainerOption {
	return func(container *ContainerManager) error {
		// Si no estamos en Jenkins vamos a alertar al usuario de que no está usando
		// la clave SSH que necesita el comando.
		if config.SSHAgentSocket() == "" && !config.Jenkins() {
			log.Warning("WARNING: No SSH_AUTH_SOCK defined in the environment. Start an ssh-agent to share the SSH keys with the tools.")
			return nil
		}

		if config.SSHAgentSocket() != "" {
			pid := filepath.Ext(config.SSHAgentSocket())
			inside := fmt.Sprintf("/tmp/ssh-sock/agent.%s", pid)

			container.env["SSH_AUTH_SOCK"] = inside
			container.volumes[config.SSHAgentSocket()] = inside
		}

		return nil
	}
}

func WithNetworkAlias(networkAlias string) ContainerOption {
	return func(container *ContainerManager) error {
		container.networkAlias = networkAlias
		return nil
	}
}

func WithoutTTY() ContainerOption {
	return func(container *ContainerManager) error {
		container.noTTY = true
		return nil
	}
}

func WithSharedWorkspace() ContainerOption {
	return func(container *ContainerManager) error {
		wd, err := os.Getwd()
		if err != nil {
			return errors.Trace(err)
		}

		container.volumes[wd] = "/workspace"
		container.workdir = "/workspace"

		return nil
	}
}

func WithEnv(name, value string) ContainerOption {
	return func(container *ContainerManager) error {
		container.env[name] = value
		return nil
	}
}

func WithVolume(source, inside string) ContainerOption {
	return func(container *ContainerManager) error {
		container.volumes[source] = inside
		return nil
	}
}

func WithSharedGcloud() ContainerOption {
	return func(container *ContainerManager) error {
		configPath := fmt.Sprintf("%s/.config/gcloud", config.Home())
		if ok, err := hasConfig(configPath); err != nil {
			return errors.Trace(err)
		} else if ok {
			container.volumes[configPath] = "/home/container/.config/gcloud"
		}

		configPath = fmt.Sprintf("%s/.gsutil", config.Home())
		if ok, err := hasConfig(configPath); err != nil {
			return errors.Trace(err)
		} else if ok {
			container.volumes[configPath] = "/home/container/.gsutil"
		}

		configPath = fmt.Sprintf("%s/.kube", config.Home())
		if ok, err := hasConfig(configPath); err != nil {
			return errors.Trace(err)
		} else if ok {
			container.volumes[configPath] = "/home/container/.kube"
		}

		return nil
	}
}

func WithPort(source, inside int64) ContainerOption {
	return func(container *ContainerManager) error {
		container.ports[source] = inside
		return nil
	}
}

func WithSharedGopath() ContainerOption {
	return func(container *ContainerManager) error {
		if config.ProjectPackage() != "" {
			hostBin := fmt.Sprintf("%s/.actools/cache-%s/bin", config.Home(), config.ProjectName())
			container.volumes[hostBin] = "/go/bin"
			if err := os.MkdirAll(hostBin, 0777); err != nil {
				return errors.Trace(err)
			}

			hostPkg := fmt.Sprintf("%s/.actools/cache-%s/pkg", config.Home(), config.ProjectName())
			container.volumes[hostPkg] = "/go/pkg"
			if err := os.MkdirAll(hostPkg, 0777); err != nil {
				return errors.Trace(err)
			}

			cachePkg := fmt.Sprintf("%s/.actools/cache-%s/cache", config.Home(), config.ProjectName())
			container.volumes[cachePkg] = "/home/container/.cache"
			if err := os.MkdirAll(cachePkg, 0777); err != nil {
				return errors.Trace(err)
			}
		}

		return nil
	}
}

func WithWorkdir(workdir string) ContainerOption {
	return func(container *ContainerManager) error {
		container.userWorkdir = workdir
		return nil
	}
}

func WithPersistence() ContainerOption {
	return func(container *ContainerManager) error {
		container.persistent = true
		return nil
	}
}

func hasConfig(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if !os.IsNotExist(err) {
			return false, errors.Trace(err)
		}

		return false, nil
	}

	return true, nil
}

func WithStandardHome() ContainerOption {
	return func(container *ContainerManager) error {
		container.env["HOME"] = "/home/container"
		return nil
	}
}
