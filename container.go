package daap

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
)

// Container represents a container on a machine.
type Container struct {
	Image string
	Args  Args
	container.ContainerCreateCreatedBody
}

// NewContainer requires necessary arguments to create a container.
func NewContainer(img string, args Args) *Container {
	return &Container{
		Image: img,
		Args:  args,
	}
}

// PullImage pulls specified image to this container.
func (c *Container) PullImage(ctx context.Context) (<-chan ImagePullResponsePayload, error) {
	dkclient, err := c.Args.Machine.CreateClient()
	if err != nil {
		return nil, err
	}
	defer dkclient.Close()
	rc, err := dkclient.ImagePull(ctx, c.Image, types.ImagePullOptions{})
	if err != nil {
		return nil, err
	}

	stream := make(chan ImagePullResponsePayload)
	scanner := bufio.NewScanner(rc)
	go func() {
		defer close(stream)
		defer rc.Close()
		for scanner.Scan() {
			payload := ImagePullResponsePayload{}
			json.Unmarshal(scanner.Bytes(), &payload)
			stream <- payload
		}
	}()

	return stream, nil
}

// Create creates container itself on a machine specified as args for NewContainer.
func (c *Container) Create(ctx context.Context) error {
	dkclient, err := c.Args.Machine.CreateClient()
	if err != nil {
		return err
	}
	defer dkclient.Close()
	c.ContainerCreateCreatedBody, err = dkclient.ContainerCreate(
		ctx,
		&container.Config{
			Image:        c.Image,
			Env:          []string{}, // Env variables with formatted "key=value"
			Tty:          true,       // To keep alive
			AttachStdout: true,
			AttachStderr: true,
		},
		&container.HostConfig{
			Mounts: []mount.Mount{}, // e.g. VolumesFrom
		},
		&network.NetworkingConfig{},
		c.Args.Name,
	)
	return err
}

// Start starts a created container.
func (c *Container) Start(ctx context.Context) error {
	dkclient, err := c.Args.Machine.CreateClient()
	if err != nil {
		return err
	}
	defer dkclient.Close()
	return dkclient.ContainerStart(ctx, c.ID, types.ContainerStartOptions{})
}

// Exec executes specified command on this container.
// Before calling "Exec", this container must be created and started.
func (c *Container) Exec(ctx context.Context, cmd ...string) (<-chan HijackedStreamPayload, error) {

	dkclient, err := c.Args.Machine.CreateClient()
	if err != nil {
		return nil, err
	}
	defer dkclient.Close()

	execute, err := dkclient.ContainerExecCreate(ctx, c.ID, types.ExecConfig{
		Cmd:          cmd,
		AttachStdout: true,
		AttachStderr: true,
	})
	if err != nil {
		return nil, fmt.Errorf("Exec Create Error: %v", err)
	}

	hijacked, err := dkclient.ContainerExecAttach(ctx, execute.ID, types.ExecConfig{})
	if err != nil {
		return nil, fmt.Errorf("Exec Attach Error: %v", err)
	}
	// defer hijacked.Close()

	err = dkclient.ContainerExecStart(ctx, execute.ID, types.ExecStartCheck{})
	if err != nil {
		hijacked.Close()
		return nil, fmt.Errorf("Exec Start Error: %v", err)
	}

	stream := make(chan HijackedStreamPayload)
	scanner := bufio.NewScanner(hijacked.Reader)
	go func() {
		defer close(stream)
		defer hijacked.Close()
		for scanner.Scan() {
			stream <- CreatePayloadFromRawBytes(MIXED, scanner.Bytes())
		}
		return
	}()
	return stream, nil
}
