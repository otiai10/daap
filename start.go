package daap

import (
	"context"

	"github.com/docker/docker/api/types"
)

// Start starts a created container.
func (c *Container) Start(ctx context.Context) error {
	dkclient, err := c.getClient()
	if err != nil {
		return err
	}
	defer dkclient.Close()
	return dkclient.ContainerStart(ctx, c.ID, types.ContainerStartOptions{})
}
