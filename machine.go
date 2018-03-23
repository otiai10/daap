package daap

import (
	"os"
)

// Machine represents a machine config holder.
type Machine interface {
	Host() string
	CertPath() string
	Version() string
}

// EnvMachine ...
type EnvMachine struct{}

// Host ...
func (m EnvMachine) Host() string {
	return os.Getenv("DOCKER_HOST")
}

// CertPath ...
func (m EnvMachine) CertPath() string {
	return os.Getenv("DOCKER_CERT_PATH")
}

// Version ...
func (m EnvMachine) Version() string {
	return ""
}
