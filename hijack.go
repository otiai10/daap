package daap

// HijackedStreamType ...
type HijackedStreamType uint8

const (
	// STDIN represents stdin: &0
	STDIN HijackedStreamType = 0
	// STDOUT represents stdout: &1
	STDOUT HijackedStreamType = 1
	// STDERR represents stderr: &2
	STDERR HijackedStreamType = 2
	// MIXED represents mixed io both of stdout and stderr
	MIXED HijackedStreamType = STDOUT | STDERR
)

// HijackedStreamPayload ...
type HijackedStreamPayload struct {
	Type HijackedStreamType
	Data []byte
}

// CreatePayloadFromRawBytes ...
// https://docs.docker.com/engine/api/v1.30/#operation/ContainerAttach
func CreatePayloadFromRawBytes(defaultType HijackedStreamType, raw []byte) HijackedStreamPayload {
	if len(raw) >= 8 {
		header := [4]byte{}
		copy(header[:], raw[:4])
		switch header {
		case [4]byte{uint8(STDOUT), 0, 0, 0}:
			return HijackedStreamPayload{Type: STDOUT, Data: raw[8:]}
		case [4]byte{uint8(STDERR), 0, 0, 0}:
			return HijackedStreamPayload{Type: STDERR, Data: raw[8:]}
		case [4]byte{uint8(STDIN), 0, 0, 0}:
			return HijackedStreamPayload{Type: STDIN, Data: raw[8:]}
		}
	}
	return HijackedStreamPayload{
		Type: defaultType,
		Data: raw,
	}
}
