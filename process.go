package daap

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

// Process represents a docker container dealt as a process.
//
// A Process can't be canceled after calling its `Run`, `Output` methods.
type Process struct {
	// Image is a container image of this process
	Image       string
	Args        Args
	Log         io.ReadWriter
	Stdout      io.ReadWriter
	Stderr      io.ReadWriter
	hijackedOut types.HijackedResponse
	hijackedErr types.HijackedResponse
	client      *client.Client
	ID          string
	Remove      bool
}

// Args represents argument for the process, representing machine (where to run), input (what to mount), output (where to output).
type Args struct {
	Machine *MachineConfig
	Env     []string
	Mounts  []Mount
	Name    string
}

// NewProcess returns the Process struct to execute the named image with the given arguments.
//
// It sets only the Image and Args in the returned structure.
func NewProcess(img string, args Args) *Process {
	return &Process{
		Image:  img,
		Args:   args,
		Remove: true,
	}
}

// Output runs the process and returns its standard output.
func (p *Process) Output() ([]byte, error) {
	return []byte{}, nil
}

// Run starts the specified container, and waits for it to complete.
//
// The returned error is nil if the process runs, has no problems copying stdout and stderr, and exits with a zero exit status.
func (p *Process) Run(ctx context.Context) error {
	if err := p.Start(ctx); err != nil {
		return err
	}
	if err := p.Wait(ctx); err != nil {
		return err
	}
	return nil
}

// Start starts the specified container but does not wait for it to complete.
//
// The `Wait` method will return the exit code and release associated resources once the command exits.
func (p *Process) Start(ctx context.Context) error {

	// Validation
	if p.Args.Machine == nil {
		return fmt.Errorf("Machine is not defined")
	}

	// Crete Docker Client
	client, err := p.Args.Machine.CreateClient()
	if err != nil {
		return fmt.Errorf("failed to create client: %s", err.Error())
	}
	p.client = client
	p.Stdout = new(bytes.Buffer)
	p.Stderr = new(bytes.Buffer)
	p.Log = new(bytes.Buffer)

	// Ensure Image existing on the host
	r, err := client.ImagePull(ctx, p.Image, types.ImagePullOptions{})
	if err != nil {
		return fmt.Errorf("failed to pull image on the host: %s", err.Error())
	}
	defer r.Close()
	for scanner := bufio.NewScanner(r); scanner.Scan(); {
		p.Log.Write(append(scanner.Bytes(), []byte("\n")...))
	}

	// Create container
	if p.Args.Name == "" {
		p.Args.Name = fmt.Sprintf("%s_%d", strings.Replace(p.Image, "/", "_", -1), time.Now().Unix())
	}
	cntnr, err := client.ContainerCreate(
		ctx,
		p.containerConfig(),
		p.hostConfig(),
		p.networkConfig(),
		p.Args.Name,
	)
	if err != nil {
		return fmt.Errorf("failed to create container: %s", err.Error())
	}
	p.ID = cntnr.ID

	hijackedOut, err := client.ContainerAttach(ctx, cntnr.ID, types.ContainerAttachOptions{Stream: true, Stdout: true})
	if err != nil {
		return fmt.Errorf("failed to attach container STDOUT: %s", err.Error())
	}
	p.hijackedOut = hijackedOut

	hijackedErr, err := client.ContainerAttach(ctx, cntnr.ID, types.ContainerAttachOptions{Stream: true, Stderr: true})
	if err != nil {
		return fmt.Errorf("failed to attach container STDERR: %s", err.Error())
	}
	p.hijackedErr = hijackedErr

	if err := client.ContainerStart(ctx, cntnr.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	return nil
}

// StderrPipe returns a pipe that will be connected to the container's standard error when the command starts.
func (p *Process) StderrPipe() io.Reader {
	return p.hijackedErr.Reader
}

// StdoutPipe returns a pipe that will be connected to the container's standard output when the command starts.
//
// Wait will close the pipe after seeing the container exit, so most callers need not close the pipe themselves; however, an implication is that it is incorrect to call Wait before all reads from the pipe have completed. For the same reason, it is incorrect to call Run when using StdoutPipe. See the example for idiomatic usage.
func (p *Process) StdoutPipe() io.Reader {
	return p.hijackedOut.Reader
}

// Wait waits for the container to exit. It must have been started by `Start`.
//
// The returned error is nil if the command runs, has no problems copying stdin, stdout, and stderr, and exits with a zero exit status.
//
// If the command fails to run or doesn't complete successfully, the error is of type *ExitError. Other error types may be returned for I/O problems.
//
// If c.Stdin is not an *os.File, Wait also waits for the I/O loop copying from c.Stdin into the process's standard input to complete.
//
// Wait releases any resources associated with the Cmd.
func (p *Process) Wait(ctx context.Context) error {
	p.drain(p.hijackedOut, p.Stdout)
	p.drain(p.hijackedErr, p.Stderr)
	return p.cleanup(ctx)
}

// --- private ---
func (p *Process) containerConfig() *container.Config {
	return &container.Config{
		Image: p.Image,
		Env:   p.Args.Env,
	}
}

func (p *Process) hostConfig() *container.HostConfig {
	mounts := []mount.Mount{}
	for _, m := range p.Args.Mounts {
		mounts = append(mounts, m.ToDockerAPITypeMount())
	}
	return &container.HostConfig{
		Mounts: mounts,
	}
}

func (p *Process) networkConfig() *network.NetworkingConfig {
	return &network.NetworkingConfig{}
}

func (p *Process) drain(hijacked types.HijackedResponse, dest io.Writer) {
	defer hijacked.Close()
	for scanner := bufio.NewScanner(hijacked.Reader); scanner.Scan(); {
		b := append(scanner.Bytes(), []byte("\n")...)
		// Hijacked stream has a 8-bytes-length header
		// See https://docs.docker.com/engine/api/v1.30/#operation/ContainerAttach for more information
		if len(b) <= 8 {
			dest.Write(b)
			continue
		}
		var fixedHeader [4]byte
		copy(fixedHeader[:], b[:4])
		switch fixedHeader {
		case [4]byte{0, 0, 0, 0}, [4]byte{1, 0, 0, 0}, [4]byte{2, 0, 0, 0}:
			dest.Write(b[8:])
		default:
			dest.Write(b)
		}
	}
}

// "cleanup" cleans up everything on container layer,
// it doesn't care about hijacked stdout/stderr layers.
func (p *Process) cleanup(ctx context.Context) error {
	res, err := p.client.ContainerInspect(ctx, p.ID)
	if err != nil {
		return err
	}
	if res.ContainerJSONBase.State.Running {
		return fmt.Errorf("This container is still running, not to be cleanup")
	}
	if err := p.client.ContainerRemove(ctx, p.ID, types.ContainerRemoveOptions{}); err != nil {
		return err
	}

	// if p.Remove {
	if _, err := p.client.ImageRemove(ctx, p.Image, types.ImageRemoveOptions{
		Force:         true,
		PruneChildren: true,
	}); err != nil {
		return err
	}
	// }

	return nil
}
