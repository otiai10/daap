package daap

import (
	"bufio"
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
)

// Execution ...
type Execution struct {
	Inline string
	Script string
	Env    []string
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
		Env:          execution.Env,
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
