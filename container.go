package daap

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"

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
func (c *Container) PullImage(ctx context.Context, out io.Writer) error {
	dkclient, err := c.Args.Machine.CreateClient()
	if err != nil {
		return err
	}
	defer dkclient.Close()
	rc, err := dkclient.ImagePull(ctx, c.Image, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	defer rc.Close()

	if out == nil {
		out = bytes.NewBuffer(nil)
	}
	scanner := bufio.NewScanner(rc)
	for scanner.Scan() {
		fmt.Fprintf(out, "\r%s", scanner.Text())
	}
	fmt.Fprintf(out, "\n")

	return scanner.Err()
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
func (c *Container) Exec(ctx context.Context, out chan<- []byte, cmd ...string) error {

	dkclient, err := c.Args.Machine.CreateClient()
	if err != nil {
		return err
	}
	defer dkclient.Close()

	execute, err := dkclient.ContainerExecCreate(ctx, c.ID, types.ExecConfig{
		Cmd:          cmd,
		AttachStdout: true,
		AttachStderr: true,
	})
	if err != nil {
		return fmt.Errorf("Exec Create Error: %v", err)
	}

	hijacked, err := dkclient.ContainerExecAttach(ctx, execute.ID, types.ExecConfig{})
	if err != nil {
		return fmt.Errorf("Exec Attach Error: %v", err)
	}
	// defer hijacked.Close()

	err = dkclient.ContainerExecStart(ctx, execute.ID, types.ExecStartCheck{})
	if err != nil {
		hijacked.Close()
		return fmt.Errorf("Exec Start Error: %v", err)
	}

	scanner := bufio.NewScanner(hijacked.Reader)

	// If out is not given, it executes syncly.
	if out == nil {
		for scanner.Scan() {
			// TODO: Pool output somewhere
			fmt.Println(scanner.Text())
		}
		hijacked.Close()
		return scanner.Err()
	}

	// If out is given, delegate output to the caller.
	go func() {
		for scanner.Scan() {
			out <- scanner.Bytes()
		}
		hijacked.Close()
		close(out)
	}()

	return nil
}
