package daap

import "github.com/docker/docker/api/types/mount"

// Volume represents --volume source:target,
// shorthand function to create "types/mount".Mount
func Volume(source, target string, readonly ...bool) mount.Mount {
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
