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

	// mounts := []mount.Mount{}
	// for _, m := range c.Args.Mounts {
	// 	mounts = append(mounts, m.ToDockerAPITypeMount())
	// }
	// c.ContainerCreateCreatedBody, err = dkclient.ContainerCreate(
	// 	ctx,
	// 	&container.Config{
	// 		Image:        c.Image,
	// 		Env:          []string{}, // Env variables with formatted "key=value"
	// 		Tty:          true,       // To keep alive
	// 		AttachStdout: true,
	// 		AttachStderr: true,
	// 	},
	// 	&container.HostConfig{
	// 		Mounts:      mounts,
	// 		NetworkMode: "host", // TODO: make configurable
	// 		Privileged:  true,
	// 	},
	// 	&network.NetworkingConfig{},
	// 	c.Args.Name,
	// )

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
		// &container.Config{
		// 	Image:        c.Image,
		// 	Env:          []string{}, // Env variables with formatted "key=value"
		// 	Tty:          true,       // To keep alive
		// 	AttachStdout: true,
		// 	AttachStderr: true,
		// },
		config.Host,
		// &container.HostConfig{
		// 	Mounts:      mounts,
		// 	NetworkMode: "host", // TODO: make configurable
		// 	Privileged:  true,
		// },
		config.Network,
		config.Name,
	)

	return err
}
