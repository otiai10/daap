package daap

import (
	"archive/tar"
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

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

// Upload uploads local file to this container.
func (c *Container) Upload(ctx context.Context, src *os.File, destdir string) error {
	dkclient, err := c.Args.Machine.CreateClient()
	if err != nil {
		return err
	}
	defer dkclient.Close()

	stat, err := src.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file stat of source file: %v", err)
	}

	buf := bytes.NewBuffer(nil)
	tarwriter := tar.NewWriter(buf)
	header := &tar.Header{
		Name: stat.Name(),
		Mode: int64(stat.Mode()),
		Size: stat.Size(),
	}

	if err := tarwriter.WriteHeader(header); err != nil {
		return fmt.Errorf("failed to write header of tar: %v", err)
	}

	if _, err := io.Copy(tarwriter, src); err != nil {
		return fmt.Errorf("failed to write content of source file as tar: %v", err)
	}

	if err := dkclient.CopyToContainer(ctx, c.ID, destdir, buf, types.CopyToContainerOptions{}); err != nil {
		return fmt.Errorf("failed to copy file to the container: %v", err)
	}

	return nil
}

// genExecCommand ...
func (c *Container) genExecCommand(ctx context.Context, execution Execution) ([]string, error) {

	if execution.Inline == "" && execution.Script == "" {
		return nil, fmt.Errorf("either of `inline` or `script` must be specified as an execution")
	}
	if execution.Inline != "" {
		return []string{"bash", "-c", execution.Inline}, nil
	}

	script, err := os.Open(execution.Script)
	if err != nil {
		return nil, fmt.Errorf("failed to open your script file: %v", err)
	}
	defer script.Close()

	if err := c.Upload(ctx, script, "/"); err != nil {
		return nil, fmt.Errorf("failed to upload: %v", err)
	}

	// TODO: Fix this hard coding of using "sh"
	//       It might be determined by extension of filename.
	cmd := []string{"sh", "/" + filepath.Base(script.Name())}
	return cmd, nil
}

// Exec executes specified command on this container.
// Before calling "Exec", this container must be created and started.
func (c *Container) Exec(ctx context.Context, execution Execution) (<-chan HijackedStreamPayload, error) {

	dkclient, err := c.Args.Machine.CreateClient()
	if err != nil {
		return nil, err
	}
	defer dkclient.Close()

	cmd, err := c.genExecCommand(ctx, execution)
	if err != nil {
		return nil, err
	}

	execute, err := dkclient.ContainerExecCreate(ctx, c.ID, types.ExecConfig{
		Cmd:          cmd,
		AttachStdout: true,
		AttachStderr: true,
	})
	if err != nil {
		return nil, fmt.Errorf("Exec Create Error: %v", err)
	}

	hijacked, err := dkclient.ContainerExecAttach(ctx, execute.ID, types.ExecStartCheck{})
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
		payload := HijackedStreamPayload{Type: MIXED}
		for scanner.Scan() {
			// If raw bytes doesn't have header, use previous io type as a default.
			payload = CreatePayloadFromRawBytes(payload.Type, scanner.Bytes())
			stream <- payload
		}
		return
	}()
	return stream, nil
}
