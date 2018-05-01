package daap

import (
	"net/http"
	"path/filepath"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/tlsconfig"
)

// Container represents a container on a machine.
type Container struct {
	Image string
	// Args    Args
	Machine Machine
	container.ContainerCreateCreatedBody

	// RetryCount for Exec (ExecCreate, ExecAttach).
	// See "retry.go" for more information.
	RetryCount int
}

// NewContainer creates a definition of a container.
// It doesn't create an actual container yet.
// Call *Container.Create to create one.
func NewContainer(img string, machine Machine) *Container {
	return &Container{
		Image:   img,
		Machine: machine,
	}
}

// getClient ...
func (c *Container) getClient() (*client.Client, error) {
	tlsc, err := tlsconfig.Client(tlsconfig.Options{
		CAFile:             c.pathTo("ca.pem"),
		CertFile:           c.pathTo("cert.pem"),
		KeyFile:            c.pathTo("key.pem"),
		InsecureSkipVerify: false,
	})
	if err != nil {
		return nil, err
	}
	httpclient := &http.Client{Transport: &http.Transport{TLSClientConfig: tlsc}}
	headers := map[string]string{}
	return client.NewClient(c.Machine.Host(), c.Machine.Version(), httpclient, headers)
}

// pathTo ...
func (c *Container) pathTo(fname string) string {
	return filepath.Join(c.Machine.CertPath(), fname)
}
