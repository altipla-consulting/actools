package docker

type ContainerOption func(container *ContainerManager)

func WithNetwork(networkName string) ContainerOption {
	return func(container *ContainerManager) {
		container.network = Network(networkName)
	}
}
