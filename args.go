package daap

//
// // Process represents a docker container dealt as a process.
// //
// // A Process can't be canceled after calling its `Run`, `Output` methods.
// type Process struct {
// 	// Image is a container image of this process
// 	Image       string
// 	Args        Args
// 	Log         io.ReadWriter
// 	Stdout      io.ReadWriter
// 	Stderr      io.ReadWriter
// 	hijackedOut types.HijackedResponse
// 	hijackedErr types.HijackedResponse
// 	client      *client.Client
// 	finished    chan error
// 	ID          string
// 	Remove      bool
// }

// Args represents argument for the process, representing machine (where to run), input (what to mount), output (where to output).
type Args struct {
	Machine *MachineConfig
	Env     []string
	Mounts  []Mount
	Name    string
}
