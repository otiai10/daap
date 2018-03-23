package daap

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
)

// Upload uploads local file to this container.
func (c *Container) Upload(ctx context.Context, src *os.File, destdir string) error {
	dkclient, err := c.getClient()
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
