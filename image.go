package daap

import (
	"bufio"
	"context"
	"encoding/json"

	"github.com/docker/docker/api/types"
)

// PullImage pulls specified image to this container.
func (c *Container) PullImage(ctx context.Context) (<-chan ImagePullResponsePayload, error) {
	dkclient, err := c.getClient()
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
