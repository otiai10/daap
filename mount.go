package daap

import "github.com/docker/docker/api/types/mount"

// MountType == types.mount.Type
type MountType string

const (
	// TypeBind is the type for mounting host dir
	TypeBind MountType = "bind"
	// TypeVolume is the type for remote storage volumes
	TypeVolume MountType = "volume"
	// TypeTmpfs is the type for mounting tmpfs
	TypeTmpfs MountType = "tmpfs"
)

// Mount ...
type Mount struct {
	Type   MountType
	Source string
	Target string
	// ReadOnly    bool        `json:",omitempty"`
	// Consistency Consistency `json:",omitempty"`
	// BindOptions   *BindOptions   `json:",omitempty"`
	// VolumeOptions *VolumeOptions `json:",omitempty"`
	// TmpfsOptions  *TmpfsOptions  `json:",omitempty"`
}

// ToDockerAPITypeMount converts daap.Mount to types.mount.Mount
func (m Mount) ToDockerAPITypeMount() mount.Mount {
	return mount.Mount{
		Type:   mount.Type(m.Type),
		Source: m.Source,
		Target: m.Target,
	}
}
