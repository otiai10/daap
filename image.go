package daap

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/docker/docker/api/types"
)

// PullImage pulls specified image to this container.
func (c *Container) PullImage(ctx context.Context) (<-chan ImagePullResponsePayload, error) {
	dkclient, err := c.getClient()
	if err != nil {
		return nil, err
	}
	defer dkclient.Close()

	var progress io.ReadCloser
	err = c.retry(func() error {
		resp, err := dkclient.ImagePull(ctx, c.Image, types.ImagePullOptions{})
		if err != nil {
			return err
		}
		progress = resp
		return nil
	}, 0, nil)
	if err != nil {
		return nil, fmt.Errorf("Image Pull Error: %v", err)
	}

	stream := make(chan ImagePullResponsePayload)
	scanner := bufio.NewScanner(progress)
	go func() {
		defer close(stream)
		defer progress.Close()
		for scanner.Scan() {
			payload := ImagePullResponsePayload{}
			json.Unmarshal(scanner.Bytes(), &payload)
			stream <- payload
		}
	}()

	return stream, nil
}
