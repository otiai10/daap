package daap

import "github.com/docker/docker/api/types/mount"

// MountVolume represents --volume source:target,
// shorthand function to create "types/mount".Mount
func MountVolume(source, target string, readonly ...bool) mount.Mount {
	readonly = append(readonly, false)
	return mount.Mount{
		Type:     mount.TypeBind,
		Source:   source,
		Target:   target,
		ReadOnly: readonly[0],
	}
}

// MountVolumesFrom represents --volumes-from containerID,
// shorthand function to create "types/mount".Mount
func MountVolumesFrom(containerID string) mount.Mount {
	return mount.Mount{
		Type:   mount.TypeVolume,
		Source: containerID,
	}
}
