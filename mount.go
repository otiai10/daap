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
	Type     MountType
	Source   string
	Target   string
	ReadOnly bool
	// Consistency Consistency `json:",omitempty"`
	// BindOptions   *BindOptions   `json:",omitempty"`
	// VolumeOptions *VolumeOptions `json:",omitempty"`
	// TmpfsOptions  *TmpfsOptions  `json:",omitempty"`
}

// ToDockerAPITypeMount converts daap.Mount to types.mount.Mount
func (m Mount) ToDockerAPITypeMount() mount.Mount {
	return mount.Mount{
		Type:     mount.Type(m.Type),
		Source:   m.Source,
		Target:   m.Target,
		ReadOnly: m.ReadOnly,
	}
}

// Volume represents --volume source:target
func Volume(source, target string, readonly ...bool) Mount {
	readonly = append(readonly, false)
	return Mount{
		Type:     TypeBind,
		Source:   source,
		Target:   target,
		ReadOnly: readonly[0],
	}
}

// VolumesFrom represents --volumes-from containerID
func VolumesFrom(containerID string) Mount {
	return Mount{
		Type:   TypeVolume,
		Source: containerID,
	}
}
