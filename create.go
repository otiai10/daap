package daap

import (
	"context"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
)

// CreateConfig ...
type CreateConfig struct {
	Name      string
	Container *container.Config
	Host      *container.HostConfig
	Network   *network.NetworkingConfig
}

// Create creates container itself on a machine specified as args for NewContainer.
func (c *Container) Create(ctx context.Context, config CreateConfig) error {

	dkclient, err := c.getClient()
	if err != nil {
		return err
	}
	defer dkclient.Close()

	if config.Container == nil {
		config.Container = &container.Config{}
	}
	if config.Host == nil {
		config.Host = &container.HostConfig{}
	}
	if config.Network == nil {
		config.Network = &network.NetworkingConfig{}
	}

	// Necessary Overwrite
	config.Container.Image = c.Image

	// TEMP: Overwrite host config for daap
	config.Container.Tty = true
	config.Container.AttachStdout = true
	config.Container.AttachStderr = true

	// TEMP: Overwrite host config for daap
	config.Host.NetworkMode = "host"
	config.Host.Privileged = true

	c.ContainerCreateCreatedBody, err = dkclient.ContainerCreate(
		ctx,
		config.Container,
		config.Host,
		config.Network,
		config.Name,
	)

	return err
}
