package daap

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

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
	return c.stream(hijacked)
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

// stream drains hijacked stdout/stderr response.
func (c *Container) stream(hijacked types.HijackedResponse) (<-chan HijackedStreamPayload, error) {
	stream := make(chan HijackedStreamPayload)
	go func() {
		buf := bufio.NewReader(hijacked.Reader)
		payload := HijackedStreamPayload{Type: MIXED}
		for {
			b, err := buf.ReadBytes('\n')
			// If raw bytes doesn't have header, use previous io type as a default.
			if err == nil {
				payload = CreatePayloadFromRawBytes(payload.Type, b)
				stream <- payload
			} else {
				if err != io.EOF {
					log.Println("Buffer Error:", err)
				}
				break
			}
		}
		hijacked.Close()
		close(stream)
		return
	}()
	return stream, nil
}
