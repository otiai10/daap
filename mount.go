package daap

import "github.com/docker/docker/api/types/mount"

// Bind represents --volume source:target,
// shorthand function to create "types/mount".Mount
func Bind(source, target string, readonly ...bool) mount.Mount {
	readonly = append(readonly, false)
	return mount.Mount{
		Type:     mount.TypeBind,
		Source:   source,
		Target:   target,
		ReadOnly: readonly[0],
	}
}

// VolumesFrom represents --volumes-from containerID,
// shorthand function to create "types/mount".Mount
func VolumesFrom(containerID string) mount.Mount {
	return mount.Mount{
		Type:   mount.TypeVolume,
		Source: containerID,
	}
}

// VolumeByName ...
func VolumeByName(volumename, target string, readonly ...bool) mount.Mount {
	readonly = append(readonly, false)
	return mount.Mount{
		Type:     mount.TypeVolume,
		Source:   volumename,
		Target:   target,
		ReadOnly: readonly[0],
	}
}
