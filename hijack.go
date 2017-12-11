package daap

// HijackedIOType ...
type HijackedIOType int

const (
	// STDOUT represents stdout: &1
	STDOUT HijackedIOType = 1
	// STDERR represents stderr: &2
	STDERR HijackedIOType = 2
	// MIXED represents mixed io both of stdout and stderr
	MIXED HijackedIOType = STDOUT | STDERR
)

// HijackedPayload ...
type HijackedPayload struct {
	Type HijackedIOType
	Data []byte
}
