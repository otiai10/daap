package daap

import (
	"context"
	"net/http"
	"path/filepath"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/tlsconfig"
)

// Volume ...
type Volume struct {
	Machine Machine
	Config  volume.VolumesCreateBody
	types.Volume
}

// Create ...
func (v *Volume) Create(ctx context.Context) error {
	m := v.Machine
	tlsc, err := tlsconfig.Client(tlsconfig.Options{
		CAFile:             filepath.Join(m.CertPath(), "ca.pem"),
		CertFile:           filepath.Join(m.CertPath(), "cert.pem"),
		KeyFile:            filepath.Join(m.CertPath(), "key.pem"),
		InsecureSkipVerify: false,
	})
	if err != nil {
		return err
	}
	httpclient := &http.Client{Transport: &http.Transport{TLSClientConfig: tlsc}}
	headers := map[string]string{}
	c, err := client.NewClient(m.Host(), m.Version(), httpclient, headers)
	if err != nil {
		return err
	}

	vol, err := c.VolumeCreate(ctx, v.Config)
	if err != nil {
		return err
	}

	v.Volume = vol
	return nil
}
