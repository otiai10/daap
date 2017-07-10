package daap

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/docker/docker/client"
	"github.com/docker/go-connections/tlsconfig"
)

// MachineConfig ...
type MachineConfig struct {
	Host     string
	CertPath string
	Version  string
}

// NewEnvMachine ...
func NewEnvMachine() *MachineConfig {
	return &MachineConfig{
		Host:     os.Getenv("DOCKER_HOST"),
		CertPath: os.Getenv("DOCKER_CERT_PATH"),
		Version:  "",
	}
}

// CreateClient ...
func (mc *MachineConfig) CreateClient() (*client.Client, error) {
	tlsc, err := tlsconfig.Client(tlsconfig.Options{
		CAFile:             filepath.Join(mc.CertPath, "ca.pem"),
		CertFile:           filepath.Join(mc.CertPath, "cert.pem"),
		KeyFile:            filepath.Join(mc.CertPath, "key.pem"),
		InsecureSkipVerify: false,
	})
	if err != nil {
		return nil, err
	}
	cl := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsc,
		},
	}
	headers := map[string]string{}
	return client.NewClient(mc.Host, mc.Version, cl, headers)
}
